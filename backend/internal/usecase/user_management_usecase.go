package usecase

import (
	"context"
	"errors"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserManagementUsecase interface {
	GetPuskesmasAccounts(ctx context.Context) ([]domain.User, error)
	CreatePuskesmasAccount(ctx context.Context, actorID uuid.UUID, name, email, contact, plainPassword, roleName string, puskesmasID uuid.UUID) error
	UpdatePuskesmasAccount(ctx context.Context, actorID uuid.UUID, userID uuid.UUID, name, email, contact string, activeStatus bool, roleName string, puskesmasID *uuid.UUID, password string) error
	DeletePuskesmasAccount(ctx context.Context, actorID uuid.UUID, userID uuid.UUID) error

	GetDinkesAccounts(ctx context.Context) ([]domain.User, error)
	CreateDinkesAccount(ctx context.Context, actorID uuid.UUID, name, email, contact, plainPassword string) error
	UpdateDinkesAccount(ctx context.Context, actorID uuid.UUID, userID uuid.UUID, name, email, contact string, activeStatus bool, password string) error
	DeleteDinkesAccount(ctx context.Context, actorID uuid.UUID, userID uuid.UUID) error

	GetTimFoggingAccounts(ctx context.Context) ([]domain.User, error)
}

type userManagementUsecase struct {
	userRepo  repository.UserRepository
	auditRepo repository.AuditRepository
}

func NewUserManagementUsecase(userRepo repository.UserRepository, auditRepo repository.AuditRepository) UserManagementUsecase {
	return &userManagementUsecase{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

func (u *userManagementUsecase) GetPuskesmasAccounts(ctx context.Context) ([]domain.User, error) {
	// Return both Petugas and Kepala Puskesmas
	petugas, err1 := u.userRepo.GetUsersByRole(ctx, domain.RolePetugasPuskesmas)
	kepala, err2 := u.userRepo.GetUsersByRole(ctx, domain.RoleKepalaPuskesmas)
	
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	
	return append(petugas, kepala...), nil
}

func (u *userManagementUsecase) CreatePuskesmasAccount(ctx context.Context, actorID uuid.UUID, name, email, contact, plainPassword, roleName string, puskesmasID uuid.UUID) error {
	if roleName != domain.RolePetugasPuskesmas && roleName != domain.RoleKepalaPuskesmas {
		return errors.New("role tidak valid untuk akun Puskesmas")
	}

	existing, _ := u.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return errors.New("email sudah terdaftar")
	}

	role, err := u.userRepo.GetRoleByName(ctx, roleName)
	if err != nil {
		return errors.New("role tidak ditemukan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &domain.User{
		Name:         name,
		Email:        email,
		Contact:      contact,
		PasswordHash: string(hashedPassword),
		RoleID:       role.ID,
		LocationID:   &puskesmasID,
		ActiveStatus: true,
	}

	err = u.userRepo.Create(ctx, newUser)
	if err == nil {
		_ = u.auditRepo.LogAction(ctx, nil, actorID, "CREATE_USER", "USER", newUser.ID, map[string]interface{}{
			"email": email,
			"role":  roleName,
		})
	}
	return err
}

func (u *userManagementUsecase) UpdatePuskesmasAccount(ctx context.Context, actorID uuid.UUID, userID uuid.UUID, name, email, contact string, activeStatus bool, roleName string, puskesmasID *uuid.UUID, password string) error {
	user, err := u.userRepo.GetUserByID(ctx, userID.String())
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if roleName != "" {
		if roleName != domain.RolePetugasPuskesmas && roleName != domain.RoleKepalaPuskesmas {
			return errors.New("role tidak valid untuk akun Puskesmas")
		}
		role, err := u.userRepo.GetRoleByName(ctx, roleName)
		if err == nil {
			user.RoleID = role.ID
		}
	}

	user.Name = name
	user.Email = email
	user.Contact = contact
	user.ActiveStatus = activeStatus
	if puskesmasID != nil {
		user.LocationID = puskesmasID
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}

	err = u.userRepo.Update(ctx, user)
	if err == nil {
		_ = u.auditRepo.LogAction(ctx, nil, actorID, "UPDATE_USER", "USER", user.ID, map[string]interface{}{
			"email":  email,
			"status": activeStatus,
		})
	}
	return err
}

func (u *userManagementUsecase) DeletePuskesmasAccount(ctx context.Context, actorID uuid.UUID, userID uuid.UUID) error {
	_, err := u.userRepo.GetUserByID(ctx, userID.String())
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	err = u.userRepo.Delete(ctx, userID.String())
	if err == nil {
		_ = u.auditRepo.LogAction(ctx, nil, actorID, "DELETE_USER", "USER", userID, map[string]interface{}{})
	}
	return err
}

func (u *userManagementUsecase) GetDinkesAccounts(ctx context.Context) ([]domain.User, error) {
	return u.userRepo.GetUsersByRole(ctx, domain.RoleDinkes)
}

func (u *userManagementUsecase) CreateDinkesAccount(ctx context.Context, actorID uuid.UUID, name, email, contact, plainPassword string) error {
	existing, _ := u.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return errors.New("email sudah terdaftar")
	}

	role, err := u.userRepo.GetRoleByName(ctx, domain.RoleDinkes)
	if err != nil {
		return errors.New("role tidak ditemukan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &domain.User{
		Name:         name,
		Email:        email,
		Contact:      contact,
		PasswordHash: string(hashedPassword),
		RoleID:       role.ID,
		ActiveStatus: true,
	}

	err = u.userRepo.Create(ctx, newUser)
	if err == nil {
		_ = u.auditRepo.LogAction(ctx, nil, actorID, "CREATE_USER", "USER", newUser.ID, map[string]interface{}{
			"email": email,
			"role":  domain.RoleDinkes,
		})
	}
	return err
}

func (u *userManagementUsecase) UpdateDinkesAccount(ctx context.Context, actorID uuid.UUID, userID uuid.UUID, name, email, contact string, activeStatus bool, password string) error {
	user, err := u.userRepo.GetUserByID(ctx, userID.String())
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	role, _ := u.userRepo.GetRoleByName(ctx, domain.RoleDinkes)
	if role != nil && user.RoleID != role.ID {
		return errors.New("hanya dapat mengubah akun Dinkes")
	}

	user.Name = name
	user.Email = email
	user.Contact = contact
	user.ActiveStatus = activeStatus

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}

	err = u.userRepo.Update(ctx, user)
	if err == nil {
		_ = u.auditRepo.LogAction(ctx, nil, actorID, "UPDATE_USER", "USER", user.ID, map[string]interface{}{
			"email":  email,
			"status": activeStatus,
		})
	}
	return err
}

func (u *userManagementUsecase) DeleteDinkesAccount(ctx context.Context, actorID uuid.UUID, userID uuid.UUID) error {
	user, err := u.userRepo.GetUserByID(ctx, userID.String())
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	role, _ := u.userRepo.GetRoleByName(ctx, domain.RoleDinkes)
	if role != nil && user.RoleID != role.ID {
		return errors.New("hanya dapat menghapus akun Dinkes")
	}

	err = u.userRepo.Delete(ctx, userID.String())
	if err == nil {
		_ = u.auditRepo.LogAction(ctx, nil, actorID, "DELETE_USER", "USER", userID, map[string]interface{}{})
	}
	return err
}

func (u *userManagementUsecase) GetTimFoggingAccounts(ctx context.Context) ([]domain.User, error) {
	return u.userRepo.GetUsersByRole(ctx, domain.RoleTimFogging)
}
