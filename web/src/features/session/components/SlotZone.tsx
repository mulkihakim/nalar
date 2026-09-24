import { useDroppable } from '@dnd-kit/core'
import type { OptionItem, OptionType } from '../types'
import { X, CheckCircle2, AlertCircle } from 'lucide-react'

interface SlotZoneProps {
  slotType: OptionType
  chosenOption?: OptionItem
  onClear?: () => void
  disabled?: boolean
  validationStatus?: 'correct' | 'incorrect' | null
}

export function SlotZone({
  slotType,
  chosenOption,
  onClear,
  disabled,
  validationStatus,
}: SlotZoneProps) {
  const { setNodeRef, isOver } = useDroppable({
    id: `slot-${slotType}`,
    data: {
      slot: slotType,
    },
    disabled,
  })

  const isGround = slotType === 'ground'

  // Border & background classes based on state
  let stateClasses = 'border-dashed border-2 border-slate-300 bg-slate-50/70 hover:border-slate-400'
  if (isOver) {
    stateClasses = 'border-dashed border-2 border-indigo-500 bg-indigo-50/80 ring-2 ring-indigo-200'
  } else if (validationStatus === 'correct') {
    stateClasses = 'border-2 border-emerald-500 bg-emerald-50/60 ring-2 ring-emerald-200'
  } else if (validationStatus === 'incorrect') {
    stateClasses = 'border-2 border-red-500 bg-red-50/60 ring-2 ring-red-200'
  } else if (chosenOption) {
    stateClasses = 'border-2 border-indigo-400 bg-white shadow-sm'
  }

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <span
            className={`text-xs font-bold px-2.5 py-0.5 rounded-full uppercase tracking-wider ${
              isGround
                ? 'bg-blue-100 text-blue-800 border border-blue-200'
                : 'bg-purple-100 text-purple-800 border border-purple-200'
            }`}
          >
            {isGround ? 'Slot Ground (Data/Fakta)' : 'Slot Warrant (Jembatan Penalaran)'}
          </span>
          {validationStatus === 'correct' && (
            <span className="flex items-center gap-1 text-xs text-emerald-600 font-medium">
              <CheckCircle2 className="w-3.5 h-3.5" /> Benar
            </span>
          )}
          {validationStatus === 'incorrect' && (
            <span className="flex items-center gap-1 text-xs text-red-600 font-medium">
              <AlertCircle className="w-3.5 h-3.5" /> Belum tepat
            </span>
          )}
        </div>

        {chosenOption && !disabled && onClear && (
          <button
            type="button"
            onClick={onClear}
            className="text-xs text-slate-400 hover:text-red-500 transition-colors flex items-center gap-1 px-1.5 py-0.5 rounded hover:bg-red-50"
          >
            <X className="w-3.5 h-3.5" /> Lepas
          </button>
        )}
      </div>

      <div
        ref={setNodeRef}
        className={`min-h-[120px] p-4 rounded-xl flex items-center justify-center transition-all ${stateClasses}`}
      >
        {chosenOption ? (
          <div className="w-full text-slate-900 text-sm font-medium leading-relaxed">
            {chosenOption.text}
          </div>
        ) : (
          <div className="text-center text-slate-400 text-xs py-4">
            <p className="font-medium">Tarik (Drag) opsi {isGround ? 'Ground' : 'Warrant'} ke sini</p>
            <p className="text-[11px] mt-1 text-slate-400">atau klik pada kartu opsi</p>
          </div>
        )}
      </div>
    </div>
  )
}
