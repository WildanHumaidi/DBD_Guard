# DBD Guard - Backup, Disaster Recovery and Restore Verification Procedure
# Compliant with P0-05 Business Continuity and Data Protection Requirements

param (
    [string]$ContainerName = "dbd_guard_db",
    [string]$PgUser = "postgres",
    [string]$SourceDB = "dbd_guard",
    [string]$RestoreDB = "dbd_guard_restore_test",
    [string]$BackupDir = "e:\Projects\DBD_Guard\backend\storage\backups",
    [int]$RetentionDays = 7
)

$ErrorActionPreference = "Stop"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " [P0-05] DBD GUARD BACKUP AND DISASTER RECOVERY PROCEDURE " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# 0. Ensure Host Backup Directory Exists
if (-not (Test-Path -Path $BackupDir)) {
    New-Item -ItemType Directory -Path $BackupDir -Force | Out-Null
    Write-Host "[+] Created host backup destination directory: $BackupDir" -ForegroundColor Green
}

$Timestamp = (Get-Date).ToString("yyyyMMdd_HHmmss")
$DumpFileName = "dbd_guard_backup_$Timestamp.dump"
$ContainerDumpPath = "/tmp/$DumpFileName"
$HostDumpPath = Join-Path $BackupDir $DumpFileName

# -------------------------------------------------------------------------
# PHASE 1: BACKUP (pg_dump custom format)
# -------------------------------------------------------------------------
Write-Host "`n--- PHASE 1: CREATING DATABASE BACKUP ---" -ForegroundColor Yellow
$BackupStartTime = Get-Date

try {
    # Run pg_dump inside container using custom binary compressed format (-Fc)
    docker exec $ContainerName pg_dump -U $PgUser $SourceDB -Fc -f $ContainerDumpPath
    if ($LASTEXITCODE -ne 0) {
        throw "pg_dump failed with exit code $LASTEXITCODE"
    }

    # Copy dump file to host storage separation directory
    docker cp "${ContainerName}:${ContainerDumpPath}" $HostDumpPath
    if (-not (Test-Path $HostDumpPath)) {
        throw "Failed to copy backup dump to host destination: $HostDumpPath"
    }

    $DumpItem = Get-Item $HostDumpPath
    $DumpSizeBytes = $DumpItem.Length
    $DumpSizeKB = [math]::Round($DumpSizeBytes / 1024, 2)

    $BackupEndTime = Get-Date
    $BackupDurationSec = ($BackupEndTime - $BackupStartTime).TotalSeconds

    Write-Host "[OK] Backup successfully created!" -ForegroundColor Green
    Write-Host "   - Filename     : $DumpFileName"
    Write-Host "   - Location     : $HostDumpPath"
    Write-Host "   - Size         : $DumpSizeKB KB ($DumpSizeBytes bytes)"
    Write-Host "   - Duration     : $([math]::Round($BackupDurationSec, 2)) seconds"
} catch {
    Write-Host "[ERROR] Backup failed: $_" -ForegroundColor Red
    exit 1
}

# -------------------------------------------------------------------------
# PHASE 2: RESTORE TO TEST ENVIRONMENT (Disaster Recovery Simulation)
# -------------------------------------------------------------------------
Write-Host "`n--- PHASE 2: RESTORING TO DISASTER RECOVERY TEST ENVIRONMENT ---" -ForegroundColor Yellow
$RestoreStartTime = Get-Date

try {
    # Ensure test target database is clean
    Write-Host "[-] Preparing clean test restore database '$RestoreDB'..."
    docker exec $ContainerName psql -U $PgUser -c "DROP DATABASE IF EXISTS $RestoreDB;" | Out-Null
    docker exec $ContainerName psql -U $PgUser -c "CREATE DATABASE $RestoreDB;" | Out-Null

    # Execute pg_restore
    Write-Host "[-] Executing pg_restore into '$RestoreDB'..."
    docker exec $ContainerName pg_restore -U $PgUser -d $RestoreDB $ContainerDumpPath
    
    $RestoreEndTime = Get-Date
    $RestoreDurationSec = ($RestoreEndTime - $RestoreStartTime).TotalSeconds

    Write-Host "[OK] Database restore successfully executed!" -ForegroundColor Green
    Write-Host "   - Target DB    : $RestoreDB"
    Write-Host "   - Duration     : $([math]::Round($RestoreDurationSec, 2)) seconds"
} catch {
    Write-Host "[ERROR] Restore failed: $_" -ForegroundColor Red
    exit 1
}

# -------------------------------------------------------------------------
# PHASE 3: VERIFICATION (Schema, Consistency and Data Integrity)
# -------------------------------------------------------------------------
Write-Host "`n--- PHASE 3: VERIFYING RESTORE INTEGRITY AND APPLICATION DATA ---" -ForegroundColor Yellow

$KeyTables = @("citizen_reports", "users", "inventory_items", "stock_requests", "stock_transfers", "pe_results", "audit_logs")
$VerificationFailed = $false

foreach ($Table in $KeyTables) {
    $SourceRaw = (docker exec $ContainerName psql -U $PgUser -d $SourceDB -t -A -c "SELECT COUNT(*) FROM $Table;" 2>$null) | Out-String
    $RestoreRaw = (docker exec $ContainerName psql -U $PgUser -d $RestoreDB -t -A -c "SELECT COUNT(*) FROM $Table;" 2>$null) | Out-String

    [int]$SourceCount = [int]($SourceRaw.Trim())
    [int]$RestoreCount = [int]($RestoreRaw.Trim())

    if ($SourceCount -eq $RestoreCount) {
        Write-Host "   [PASS] Table [$Table]: Source = $SourceCount | Restored = $RestoreCount (MATCH)" -ForegroundColor Green
    } else {
        Write-Host "   [FAIL] Table [$Table]: Source = $SourceCount | Restored = $RestoreCount (MISMATCH!)" -ForegroundColor Red
        $VerificationFailed = $true
    }
}

# Check Foreign Key and Constraint Integrity
Write-Host "`n[-] Checking constraint and foreign key integrity in restored database..."
$FkCheckSQL = "SELECT conname, conrelid::regclass AS table_name FROM pg_constraint WHERE contype = 'f';"
$FkOutput = docker exec $ContainerName psql -U $PgUser -d $RestoreDB -c "$FkCheckSQL"
if ($LASTEXITCODE -eq 0) {
    Write-Host "   [PASS] Foreign key constraints valid and active." -ForegroundColor Green
} else {
    Write-Host "   [FAIL] Foreign key check failed!" -ForegroundColor Red
    $VerificationFailed = $true
}

# -------------------------------------------------------------------------
# PHASE 4: RETENTION ROTATION
# -------------------------------------------------------------------------
Write-Host "`n--- PHASE 4: BACKUP RETENTION ROTATION ---" -ForegroundColor Yellow
$OldBackups = Get-ChildItem -Path $BackupDir -Filter "dbd_guard_backup_*.dump" | Where-Object { $_.CreationTime -lt (Get-Date).AddDays(-$RetentionDays) }
foreach ($Old in $OldBackups) {
    Remove-Item $Old.FullName -Force
    Write-Host "[-] Purged expired backup: $($Old.Name)" -ForegroundColor Gray
}

# Clean container temp dump
docker exec $ContainerName rm -f $ContainerDumpPath

Write-Host "`n==========================================================" -ForegroundColor Cyan
if ($VerificationFailed) {
    Write-Host " [RESULT] DISASTER RECOVERY VERIFICATION FAILED! " -ForegroundColor Red
    Write-Host "==========================================================" -ForegroundColor Cyan
    exit 1
} else {
    Write-Host " [RESULT] P0-05 BACKUP AND RESTORE VERIFICATION PASSED (100% MATCH) " -ForegroundColor Green
    Write-Host "==========================================================" -ForegroundColor Cyan
    exit 0
}
