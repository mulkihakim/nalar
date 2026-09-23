import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { BookOpen, AlertCircle, Loader2, Eye, EyeOff } from 'lucide-react'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { Label } from '../../../components/ui/label'
import { useLoginMutation } from '../api'
import { useAuth } from '../context'
import { loginSchema, type LoginFormData } from '../types'

export function LoginForm() {
  const navigate = useNavigate()
  const { login } = useAuth()
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [showPassword, setShowPassword] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      username: '',
      password: '',
    },
  })

  const loginMutation = useLoginMutation()

  const onSubmit = (data: LoginFormData) => {
    setErrorMessage(null)
    loginMutation.mutate(data, {
      onSuccess: (res) => {
        login(res.token, res.user)

        // Redirect sesuai role
        if (res.user.role === 'admin') {
          navigate({ to: '/admin/users' })
        } else if (res.user.role === 'asesor') {
          navigate({ to: '/admin/classes' })
        } else {
          navigate({ to: '/student/exams' })
        }
      },
      onError: (err) => {
        setErrorMessage(err.message || 'Login gagal. Periksa kembali username dan password Anda.')
      },
    })
  }

  return (
    <div className="w-full max-w-sm space-y-6">
      {/* Brand Header per 09-design.md §4a */}
      <div className="flex flex-col items-center space-y-2 text-center">
        <div className="w-12 h-12 rounded-xl bg-indigo-600 flex items-center justify-center text-white shadow-md shadow-indigo-100">
          <BookOpen className="w-6 h-6" />
        </div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-900">
          Nalar
        </h1>
        <p className="text-sm text-slate-500">
          Platform Latihan Argumentasi Model Toulmin
        </p>
      </div>

      {/* Card Login */}
      <div className="bg-white border border-slate-200 rounded-2xl p-6 sm:p-7 shadow-xs">
        <div className="mb-5">
          <h2 className="text-lg font-semibold text-slate-900">Masuk ke Akun</h2>
          <p className="text-xs text-slate-500 mt-0.5">
            Gunakan username dan password yang telah didaftarkan
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {errorMessage && (
            <div className="flex items-start gap-2.5 p-3 text-xs text-red-700 bg-red-50 border border-red-200 rounded-lg">
              <AlertCircle className="w-4 h-4 shrink-0 text-red-600 mt-0.5" />
              <div className="leading-snug">{errorMessage}</div>
            </div>
          )}

          <div className="space-y-1.5">
            <Label htmlFor="username" className="text-xs font-medium text-slate-700">
              Username
            </Label>
            <Input
              id="username"
              type="text"
              placeholder="Masukkan username"
              autoComplete="username"
              className="h-10 text-sm border-slate-200 focus-visible:ring-indigo-500"
              {...register('username')}
            />
            {errors.username && (
              <p className="text-xs text-red-500">{errors.username.message}</p>
            )}
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="password" className="text-xs font-medium text-slate-700">
              Password
            </Label>
            <div className="relative">
              <Input
                id="password"
                type={showPassword ? 'text' : 'password'}
                placeholder="Masukkan password"
                autoComplete="current-password"
                className="h-10 text-sm border-slate-200 focus-visible:ring-indigo-500 pr-10"
                {...register('password')}
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                tabIndex={-1}
                className="absolute right-0 top-0 h-full px-3 text-slate-400 hover:text-slate-600 flex items-center justify-center cursor-pointer transition-colors focus:outline-hidden"
                title={showPassword ? 'Sembunyikan password' : 'Lihat password'}
                aria-label={showPassword ? 'Sembunyikan password' : 'Lihat password'}
              >
                {showPassword ? (
                  <EyeOff className="w-4 h-4" />
                ) : (
                  <Eye className="w-4 h-4" />
                )}
              </button>
            </div>
            {errors.password && (
              <p className="text-xs text-red-500">{errors.password.message}</p>
            )}
          </div>

          <Button
            type="submit"
            className="w-full h-11 text-sm font-medium bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg shadow-xs transition-colors cursor-pointer mt-2"
            disabled={loginMutation.isPending}
          >
            {loginMutation.isPending ? (
              <>
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                Memproses...
              </>
            ) : (
              'Masuk'
            )}
          </Button>
        </form>
      </div>

      <div className="text-center text-xs text-slate-400">
        Nalar &copy; 2026
      </div>
    </div>
  )
}
