/**
 * Central role-based access control for the admin dashboard.
 *
 * This is the single source of truth used by:
 *  - the router guard (blocks navigation to disallowed admin pages)
 *  - the Sidebar menu (hides disallowed items)
 *
 * Access model:
 *  - admin / manager / inventory : dashboard + all pages
 *  - reception                   : dashboard + reservations (front desk) only
 *  - finance                     : the dedicated Finance page only
 *
 * Tweak the arrays below to change which roles can access which page.
 */

export type AppRole =
  | 'admin'
  | 'manager'
  | 'user'
  | 'reception'
  | 'finance'
  | 'inventory'
  | 'housekeeping'

/**
 * Allowed roles per admin route path.
 * A role not listed for an admin path is blocked from that page.
 */
export const ROUTE_ACCESS: Record<string, AppRole[]> = {
  '/admin/dashboard': ['admin', 'manager', 'reception', 'inventory', 'housekeeping'],
  '/admin/finance': ['admin', 'finance'],
  '/admin/finance/cash-bank': ['admin', 'finance'],
  '/admin/finance/general-journal': ['admin', 'finance'],
  '/admin/finance/coa': ['admin', 'finance'],
  '/admin/finance/ap-ar': ['admin', 'finance'],
  '/admin/finance/general-ledger': ['admin', 'finance'],
  '/admin/finance/profit-loss': ['admin', 'finance'],
  '/admin/finance/balance-sheet': ['admin', 'finance'],
  '/admin/reservations': ['admin', 'manager', 'reception'],
  '/admin/service-requests': ['admin', 'manager', 'reception', 'housekeeping'],
  '/admin/guestbook': ['admin', 'manager', 'reception'],
  '/admin/reservations/new': ['admin', 'manager', 'reception'],
  '/admin/rooms': ['admin', 'manager', 'housekeeping', 'reception'],
  '/admin/room-units': ['admin', 'manager', 'housekeeping', 'reception'],
  '/admin/floors': ['admin'], // Master config only for admin
  '/admin/room-types': ['admin'], // Master config only for admin
  '/admin/room-categories': ['admin'], // Master config only for admin
  '/admin/inventory': ['admin', 'manager', 'inventory'],
  '/admin/inventory-usage': ['admin', 'manager', 'inventory'],
  '/admin/promotions': ['admin', 'manager'],
  '/admin/users': ['admin', 'manager'],
  '/admin/channels': ['admin'], // OTA tech channels only for admin
  '/admin/calendar': ['admin', 'manager', 'reception'],
  '/admin/monthly-report': ['admin', 'finance'],
  '/admin/staff': ['admin', 'manager'],
  '/admin/help': ['admin', 'manager', 'inventory', 'finance', 'housekeeping'],
}

export function canAccess(role: string | undefined, path: string): boolean {
  if (!role) return false
  const allowed = ROUTE_ACCESS[path]
  if (!allowed) return true // unknown admin paths default to allowed for everyone
  return allowed.includes(role as AppRole)
}

/**
 * Cek apakah role diizinkan melakukan CUD (Create, Update, Delete).
 * Role 'manager' bersifat Read-Only / Supervisory (Monitoring & Laporan) sehingga tidak bisa memodifikasi master data.
 */
export function canMutateData(role: string | undefined): boolean {
  if (!role) return false
  return role === 'admin'
}

export function isManagerRole(role: string | undefined): boolean {
  return role === 'manager'
}

/**
 * Returns the first admin page the given role is allowed to see.
 * Used by the route guard to redirect users away from pages they can't open.
 */
export function firstAllowedPage(role: string | undefined, fallback = 'admin-dashboard'): string {
  if (!role) return fallback
  const entry = Object.entries(ROUTE_ACCESS).find(([, roles]) =>
    (roles as AppRole[]).includes(role as AppRole),
  )
  return entry ? pathToRouteName(entry[0]) : 'admin-forbidden'
}

function pathToRouteName(path: string): string {
  const names: Record<string, string> = {
    '/admin/dashboard': 'admin-dashboard',
    '/admin/finance': 'admin-finance',
    '/admin/finance/cash-bank': 'admin-finance-cash',
    '/admin/finance/general-journal': 'admin-finance-journal',
    '/admin/finance/coa': 'admin-finance-coa',
    '/admin/finance/ap-ar': 'admin-finance-apar',
    '/admin/finance/general-ledger': 'admin-finance-gl',
    '/admin/finance/profit-loss': 'admin-finance-pl',
    '/admin/finance/balance-sheet': 'admin-finance-bs',
    '/admin/reservations': 'admin-reservations',
    '/admin/guestbook': 'admin-guestbook',
    '/admin/reservations/new': 'admin-reservation-new',
    '/admin/rooms': 'admin-rooms',
    '/admin/room-types': 'admin-room-types',
    '/admin/room-categories': 'admin-room-categories',
    '/admin/inventory': 'admin-inventory',
    '/admin/promotions': 'admin-promotions',
    '/admin/users': 'admin-users',
    '/admin/channels': 'admin-channels',
    '/admin/calendar': 'admin-calendar',
    '/admin/monthly-report': 'admin-monthly-report',
    '/admin/staff': 'admin-staff',
    '/admin/help': 'admin-help',
  }
  return names[path] || 'admin-dashboard'
}