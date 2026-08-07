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

/**
 * Allowed roles per admin route path.
 * A role not listed for an admin path is blocked from that page.
 */
export const ROUTE_ACCESS: Record<string, AppRole[]> = {
  '/admin/dashboard': ['admin', 'manager', 'inventory', 'reception'],
  '/admin/finance': ['finance'],
  '/admin/reservations': ['admin', 'manager', 'inventory', 'reception'],
  '/admin/reservations/new': ['admin', 'manager', 'reception'],
  '/admin/rooms': ['admin', 'manager', 'inventory'],
  '/admin/room-types': ['admin', 'manager', 'inventory'],
  '/admin/room-categories': ['admin', 'manager', 'inventory'],
  '/admin/inventory': ['admin', 'manager', 'inventory'],
  '/admin/promotions': ['admin', 'manager', 'inventory'],
  '/admin/users': ['admin', 'manager', 'inventory'],
  '/admin/staff': ['admin', 'manager', 'inventory'],
  '/admin/help': ['admin', 'manager', 'inventory', 'finance'],
}

export function canAccess(role: string | undefined, path: string): boolean {
  if (!role) return false
  const allowed = ROUTE_ACCESS[path]
  if (!allowed) return true // unknown admin paths default to allowed for everyone
  return allowed.includes(role as AppRole)
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
    '/admin/reservations': 'admin-reservations',
    '/admin/reservations/new': 'admin-reservation-new',
    '/admin/rooms': 'admin-rooms',
    '/admin/room-types': 'admin-room-types',
    '/admin/room-categories': 'admin-room-categories',
    '/admin/inventory': 'admin-inventory',
    '/admin/promotions': 'admin-promotions',
    '/admin/users': 'admin-users',
    '/admin/staff': 'admin-staff',
    '/admin/help': 'admin-help',
  }
  return names[path] || 'admin-dashboard'
}