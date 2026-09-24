import type { ReactNode } from 'react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from './ui/dialog'

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  description?: string
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl' | '4xl'
  children: ReactNode
}

const maxWidthClasses: Record<string, string> = {
  sm: 'sm:max-w-sm',
  md: 'sm:max-w-md',
  lg: 'sm:max-w-lg',
  xl: 'sm:max-w-xl',
  '2xl': 'sm:max-w-2xl',
  '3xl': 'sm:max-w-3xl',
  '4xl': 'sm:max-w-4xl',
}

export function Modal({
  isOpen,
  onClose,
  title,
  description,
  maxWidth = 'md',
  children,
}: ModalProps) {
  const widthClass = maxWidthClasses[maxWidth] || 'sm:max-w-md'

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className={`max-h-[90vh] overflow-y-auto ${widthClass}`}>
        <DialogHeader>
          <DialogTitle className="text-base font-bold text-slate-900">{title}</DialogTitle>
          {description && (
            <DialogDescription className="text-xs text-slate-500 leading-relaxed">
              {description}
            </DialogDescription>
          )}
        </DialogHeader>
        <div className="mt-1">{children}</div>
      </DialogContent>
    </Dialog>
  )
}
