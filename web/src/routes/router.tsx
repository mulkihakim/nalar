import { createRoute, createRouter } from '@tanstack/react-router'
import { Route as rootRoute } from './__root'
import { LoginPage } from './login'
import { AdminUsersPage } from './admin/users'
import { AdminClassesPage } from './admin/classes'
import { StudentExamsPage } from './student/exams'

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: LoginPage,
})

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: LoginPage,
})

const adminUsersRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/users',
  component: AdminUsersPage,
})

const adminClassesRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/classes',
  component: AdminClassesPage,
})

const studentExamsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/student/exams',
  component: StudentExamsPage,
})

const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  adminUsersRoute,
  adminClassesRoute,
  studentExamsRoute,
])

export const router = createRouter({
  routeTree,
  defaultPreload: 'intent',
})

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}
