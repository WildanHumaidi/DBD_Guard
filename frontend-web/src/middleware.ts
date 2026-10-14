import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// Simple JWT decoder for edge runtime
function decodeJwt(token: string) {
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )
    return JSON.parse(jsonPayload)
  } catch {
    return null
  }
}

export function middleware(request: NextRequest) {
  // We use auth_token now instead of admin_token
  const token = request.cookies.get('auth_token')?.value || request.cookies.get('admin_token')?.value
  const { pathname } = request.nextUrl

  // Protected paths mapping
  const rolePaths = {
    'ROLE_ADMIN_DEVELOPER': '/developer',
    'ROLE_DINKES': '/dinkes',
    'ROLE_KADER': '/kader',
    'ROLE_PETUGAS_PUSKESMAS': '/puskesmas',
    'ROLE_KEPALA_PUSKESMAS': '/puskesmas',
    'ROLE_TIM_FOGGING': '/fogging'
  }

  const protectedPrefixes = ['/developer', '/dinkes', '/kader', '/puskesmas', '/fogging']
  const isProtectedPath = protectedPrefixes.some(prefix => pathname.startsWith(prefix))

  // 1. If accessing login while logged in, redirect to correct dashboard
  if (pathname.startsWith('/login')) {
    if (token) {
      const payload = decodeJwt(token)
      if (payload && payload.role && rolePaths[payload.role as keyof typeof rolePaths]) {
        return NextResponse.redirect(new URL(rolePaths[payload.role as keyof typeof rolePaths], request.url))
      }
    }
    return NextResponse.next()
  }

  // 2. If accessing protected paths
  if (isProtectedPath) {
    if (!token) {
      // Not logged in -> go to login
      const url = request.nextUrl.clone()
      url.pathname = '/login'
      return NextResponse.redirect(url)
    }

    // Check role authorization
    const payload = decodeJwt(token)
    if (!payload || !payload.role) {
      // Invalid token -> go to login
      const url = request.nextUrl.clone()
      url.pathname = '/login'
      return NextResponse.redirect(url)
    }

    const userRole = payload.role as keyof typeof rolePaths
    const allowedPath = rolePaths[userRole]

    if (!allowedPath) {
      // Unknown role -> public
      return NextResponse.redirect(new URL('/', request.url))
    }

    // If user is trying to access a different dashboard than their role allows
    if (!pathname.startsWith(allowedPath)) {
      return NextResponse.redirect(new URL(allowedPath, request.url))
    }
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/developer/:path*', '/dinkes/:path*', '/kader/:path*', '/puskesmas/:path*', '/fogging/:path*', '/login/:path*'],
}
