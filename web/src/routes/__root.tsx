import { useState, useEffect } from 'react'
import { createRootRoute, Outlet, Link, useNavigate, useRouterState } from '@tanstack/react-router'
import { BookOpen, Users, GraduationCap, ClipboardList, LogOut, PanelLeft, X } from 'lucide-react'
import { useAuth } from '../features/auth/context'
import { Button } from '../components/ui/button'

export const Route = createRootRoute({
  component: RootComponent,
})

function RootComponent() {
  const { user, isAuthenticated, logout } = useAuth()
  const navigate = useNavigate()
  const routerState = useRouterState()

  // State terpisah untuk desktop (collapse side-by-side) dan mobile (drawer overlay)
  const [isDesktopSidebarOpen, setIsDesktopSidebarOpen] = useState(true)
  const [isMobileDrawerOpen, setIsMobileDrawerOpen] = useState(false)

  // Tutup drawer mobile saat berpindah halaman
  useEffect(() => {
    setIsMobileDrawerOpen(false)
  }, [routerState.location.pathname])

  const isAuthPage =
    routerState.location.pathname === '/login' ||
    (!isAuthenticated && routerState.location.pathname === '/')

  const handleLogout = () => {
    logout()
    navigate({ to: '/login' })
  }

  // Sesuai 09-design.md §4a: Halaman auth standalone tanpa chrome aplikasi
  if (isAuthPage) {
    return (
      <div className="h-screen w-full flex items-center justify-center bg-slate-50 p-4 overflow-hidden">
        <Outlet />
      </div>
    )
  }

  const toggleSidebar = () => {
    if (typeof window !== 'undefined' && window.innerWidth < 768) {
      setIsMobileDrawerOpen((prev) => !prev)
    } else {
      setIsDesktopSidebarOpen((prev) => !prev)
    }
  }

  const navigationLinks = (
    <nav className="space-y-1">
      {(user?.role === 'admin' || user?.role === 'asesor') && (
        <>
          <Link
            to="/admin/users"
            className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-slate-600 hover:text-slate-900 hover:bg-slate-100 transition-colors"
            activeProps={{ className: 'bg-indigo-50 text-indigo-600 font-semibold shadow-xs' }}
          >
            <Users className="w-4 h-4 shrink-0" />
            <span>Pengguna</span>
          </Link>

          <Link
            to="/admin/classes"
            className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-slate-600 hover:text-slate-900 hover:bg-slate-100 transition-colors"
            activeProps={{ className: 'bg-indigo-50 text-indigo-600 font-semibold shadow-xs' }}
          >
            <GraduationCap className="w-4 h-4 shrink-0" />
            <span>Kelas</span>
          </Link>
        </>
      )}

      {user?.role === 'siswa' && (
        <Link
          to="/student/exams"
          className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-slate-600 hover:text-slate-900 hover:bg-slate-100 transition-colors"
          activeProps={{ className: 'bg-indigo-50 text-indigo-600 font-semibold shadow-xs' }}
        >
          <ClipboardList className="w-4 h-4 shrink-0" />
          <span>Daftar Ujian</span>
        </Link>
      )}
    </nav>
  )

  return (
    <div className="h-screen w-full flex flex-col bg-slate-50 text-slate-900 font-sans overflow-hidden">
      {/* Header Bar - Fixed 64px (h-16) per 09-design.md §4b */}
      <header className="h-16 shrink-0 border-b border-slate-200 bg-white z-30 px-3 sm:px-6 flex items-center justify-between">
        <div className="flex items-center space-x-2 sm:space-x-4">
          {isAuthenticated && (
            <button
              type="button"
              onClick={toggleSidebar}
              className="p-2 rounded-lg text-slate-500 hover:text-slate-800 hover:bg-slate-100 transition-colors cursor-pointer"
              title="Toggle Menu"
              aria-label="Toggle Sidebar"
            >
              <PanelLeft className="w-5 h-5" />
            </button>
          )}

          <Link
            to="/"
            className="flex items-center space-x-2 text-lg sm:text-xl font-bold tracking-tight text-indigo-600 hover:opacity-90 transition-opacity"
          >
            <div className="w-7 h-7 sm:w-8 sm:h-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white shadow-xs">
              <BookOpen className="w-4 h-4" />
            </div>
            <span>Nalar</span>
          </Link>
        </div>

        {/* User profile & Logout - Rapi dan responsif bahkan di layar mobile 320px */}
        <div className="flex items-center space-x-2 sm:space-x-3">
          {isAuthenticated && user && (
            <>
              <div className="flex items-center space-x-2 bg-slate-50 border border-slate-200 rounded-full py-1 px-2 sm:px-3">
                <div className="w-6 h-6 rounded-full bg-indigo-600 text-white flex items-center justify-center font-bold text-xs uppercase shadow-xs shrink-0">
                  {user.name ? user.name.charAt(0) : 'U'}
                </div>
                <span className="text-xs sm:text-sm font-semibold text-slate-800 max-w-[85px] sm:max-w-[140px] truncate">
                  {user.name}
                </span>
                <span className="hidden sm:inline-block text-[10px] font-semibold uppercase tracking-wider px-1.5 py-0.5 rounded-full bg-indigo-50 text-indigo-700 border border-indigo-200">
                  {user.role}
                </span>
              </div>

              <Button
                variant="outline"
                size="sm"
                onClick={handleLogout}
                className="cursor-pointer gap-1 px-2.5 sm:px-3 text-slate-600 hover:text-red-600 hover:bg-red-50 hover:border-red-200 transition-colors text-xs"
                title="Keluar"
              >
                <LogOut className="w-3.5 h-3.5" />
                <span className="hidden sm:inline">Keluar</span>
              </Button>
            </>
          )}
        </div>
      </header>

      {/* Body container: Sidebar (Desktop) / Drawer (Mobile) + Main Content */}
      <div className="flex-1 flex overflow-hidden relative">
        {/* MOBILE DRAWER OVERLAY (Hanya di mobile, konten utama tidak tertekan) */}
        {isAuthenticated && isMobileDrawerOpen && (
          <div
            className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-40 md:hidden transition-opacity"
            onClick={() => setIsMobileDrawerOpen(false)}
            aria-hidden="true"
          />
        )}

        {/* MOBILE DRAWER (Floating sheet di mobile, tidak menggeser konten utama) */}
        {isAuthenticated && (
          <aside
            className={`fixed inset-y-0 left-0 z-50 w-64 bg-white border-r border-slate-200 shadow-2xl flex flex-col justify-between p-4 transition-transform duration-300 ease-in-out md:hidden ${
              isMobileDrawerOpen ? 'translate-x-0' : '-translate-x-full'
            }`}
          >
            <div className="space-y-6">
              <div className="flex items-center justify-between border-b border-slate-100 pb-3">
                <div className="flex items-center space-x-2">
                  <div className="w-6 h-6 rounded-md bg-indigo-600 flex items-center justify-center text-white">
                    <BookOpen className="w-3.5 h-3.5" />
                  </div>
                  <span className="font-bold text-slate-900">Nalar Menu</span>
                </div>
                <button
                  type="button"
                  onClick={() => setIsMobileDrawerOpen(false)}
                  className="p-1 rounded-md text-slate-400 hover:text-slate-600 hover:bg-slate-100 cursor-pointer"
                  aria-label="Tutup Menu"
                >
                  <X className="w-5 h-5" />
                </button>
              </div>

              <div>
                <div className="text-[11px] font-bold text-slate-400 uppercase tracking-wider px-3 mb-2">
                  Menu Utama
                </div>
                {navigationLinks}
              </div>
            </div>

            <div className="text-[11px] text-slate-400 text-center py-2 border-t border-slate-100">
              Nalar &copy; 2026
            </div>
          </aside>
        )}

        {/* DESKTOP SIDEBAR (Side-by-side normal hanya di layar md ke atas) */}
        {isAuthenticated && (
          <aside
            className={`hidden md:flex shrink-0 bg-white border-r border-slate-200 transition-all duration-300 ease-in-out flex-col justify-between overflow-y-auto ${
              isDesktopSidebarOpen ? 'w-60 p-4' : 'w-0 p-0 border-r-0 overflow-hidden'
            }`}
          >
            <div className="space-y-6">
              <div>
                <div className="text-[11px] font-bold text-slate-400 uppercase tracking-wider px-3 mb-2">
                  Menu Utama
                </div>
                {navigationLinks}
              </div>
            </div>

            <div className="text-[11px] text-slate-400 text-center py-2 border-t border-slate-100">
              Nalar &copy; 2026
            </div>
          </aside>
        )}

        {/* Main Workspace Area: Selalu 100% full width di mobile, tidak pernah tertekan/bergeser */}
        <main className="flex-1 overflow-y-auto bg-slate-50 p-4 sm:p-6 md:p-8 flex flex-col justify-between w-full min-w-0">
          <div className="max-w-6xl mx-auto w-full">
            <Outlet />
          </div>

          <footer className="mt-12 py-4 text-center text-xs text-slate-400 border-t border-slate-200/60 max-w-6xl mx-auto w-full">
            Nalar — Platform Latihan Argumentasi Model Toulmin
          </footer>
        </main>
      </div>
    </div>
  )
}
