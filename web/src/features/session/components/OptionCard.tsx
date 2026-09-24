import { useDraggable } from '@dnd-kit/core'
import type { OptionItem } from '../types'
import { GripVertical } from 'lucide-react'

interface OptionCardProps {
  option: OptionItem
  isChosen?: boolean
  disabled?: boolean
}

export function OptionCard({ option, isChosen, disabled }: OptionCardProps) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: `opt-${option.id}`,
    data: {
      option,
    },
    disabled: disabled || isChosen,
  })

  const style = transform
    ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
        zIndex: 50,
      }
    : undefined

  const isGround = option.type === 'ground'

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className={`p-3.5 rounded-xl border transition-all text-sm select-none ${
        isChosen
          ? 'opacity-40 border-dashed border-slate-300 bg-slate-50 cursor-not-allowed'
          : isDragging
          ? 'opacity-70 border-indigo-500 bg-indigo-50 shadow-xl ring-2 ring-indigo-400 cursor-grabbing'
          : 'bg-white border-slate-200 hover:border-slate-300 hover:shadow-sm cursor-grab'
      }`}
    >
      <div className="flex items-start gap-2.5">
        <div className="text-slate-400 mt-0.5">
          <GripVertical className="w-4 h-4" />
        </div>
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <span
              className={`text-[11px] font-semibold px-2 py-0.5 rounded-full uppercase tracking-wider ${
                isGround
                  ? 'bg-blue-50 text-blue-700 border border-blue-200'
                  : 'bg-purple-50 text-purple-700 border border-purple-200'
              }`}
            >
              {isGround ? 'Ground' : 'Warrant'}
            </span>
          </div>
          <p className="text-slate-800 leading-relaxed">{option.text}</p>
        </div>
      </div>
    </div>
  )
}
