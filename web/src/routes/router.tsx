import { createRoute, createRouter } from '@tanstack/react-router'
import { Route as rootRoute } from './__root'
import { LoginPage } from './login'
import { AdminUsersPage } from './admin/users'
import { AdminClassesPage } from './admin/classes'
import { AdminMaterialsPage } from './admin/materials'
import { AdminExamsPage } from './admin/exams'
import { AdminResultsPage } from './admin/results'
import { StudentExamsPage } from './student/exams'
import { StudentExamPage } from './student/exam.$examId'

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

const adminMaterialsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/materials',
  component: AdminMaterialsPage,
})

const adminExamsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/exams',
  component: AdminExamsPage,
})

const adminResultsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/results',
  component: AdminResultsPage,
})

const studentExamsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/student/exams',
  component: StudentExamsPage,
})

const studentExamDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/student/exam/$examId',
  component: StudentExamPage,
})

const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  adminUsersRoute,
  adminClassesRoute,
  adminMaterialsRoute,
  adminExamsRoute,
  adminResultsRoute,
  studentExamsRoute,
  studentExamDetailRoute,
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
