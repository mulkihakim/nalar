import { useState } from 'react'
import { Modal } from '@/components/Modal'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ConfirmDialog } from '@/components/ConfirmDialog'
import type { SessionMode } from '../types'
import { BookOpen, Sparkles, Users } from 'lucide-react'

interface ModeSelectDialogProps {
  isOpen: boolean
  onClose: () => void
  onSelectMode: (mode: SessionMode, confirmModeChange?: boolean) => void
  isLoading?: boolean
  conflictMode?: SessionMode
  onCancelConflict?: () => void
}

export function ModeSelectDialog({
  isOpen,
  onClose,
  onSelectMode,
  isLoading,
  conflictMode,
  onCancelConflict,
}: ModeSelectDialogProps) {
  const [selectedMode, setSelectedMode] = useState<SessionMode>('standard')

  const handleStart = () => {
    onSelectMode(selectedMode, false)
  }

  return (
    <>
      <Modal
        isOpen={isOpen && !conflictMode}
        onClose={onClose}
        title="Pilih Mode Latihan"
        maxWidth="lg"
      >
        <div className="space-y-4 py-2">
          <p className="text-sm text-slate-500">
            Pilihlah mode belajar yang ingin Anda gunakan untuk mengerjakan latihan ini.
          </p>

          <div className="grid grid-cols-1 gap-3">
            {/* Mode Standar */}
            <div
              onClick={() => setSelectedMode('standard')}
              className={`flex items-start gap-4 p-4 rounded-xl border-2 cursor-pointer transition-all ${
                selectedMode === 'standard'
                  ? 'border-indigo-600 bg-indigo-50/50 shadow-sm'
                  : 'border-slate-200 hover:border-slate-300 bg-white'
              }`}
            >
              <div
                className={`p-2.5 rounded-lg ${
                  selectedMode === 'standard' ? 'bg-indigo-600 text-white' : 'bg-slate-100 text-slate-600'
                }`}
              >
                <BookOpen className="w-5 h-5" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <h4 className="font-semibold text-slate-900">Mode Standar</h4>
                  {selectedMode === 'standard' && (
                    <Badge variant="outline" className="bg-indigo-100 text-indigo-700 border-indigo-200">
                      Terpilih
                    </Badge>
                  )}
                </div>
                <p className="text-xs text-slate-500 mt-1">
                  Mengevaluasi keseluruhan argumen setelah tombol Confirm ditekan tanpa membocorkan slot mana yang salah.
                </p>
              </div>
            </div>

            {/* Mode Bantuan */}
            <div
              onClick={() => setSelectedMode('help')}
              className={`flex items-start gap-4 p-4 rounded-xl border-2 cursor-pointer transition-all ${
                selectedMode === 'help'
                  ? 'border-indigo-600 bg-indigo-50/50 shadow-sm'
                  : 'border-slate-200 hover:border-slate-300 bg-white'
              }`}
            >
              <div
                className={`p-2.5 rounded-lg ${
                  selectedMode === 'help' ? 'bg-indigo-600 text-white' : 'bg-slate-100 text-slate-600'
                }`}
              >
                <Sparkles className="w-5 h-5" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <h4 className="font-semibold text-slate-900">Mode Bantuan</h4>
                  {selectedMode === 'help' && (
                    <Badge variant="outline" className="bg-indigo-100 text-indigo-700 border-indigo-200">
                      Terpilih
                    </Badge>
                  )}
                </div>
                <p className="text-xs text-slate-500 mt-1">
                  Memberikan feedback visual pada slot ground dan warrant mana yang benar atau salah jika jawaban belum tepat.
                </p>
              </div>
            </div>

            {/* Mode Analitik Sosial */}
            <div
              onClick={() => setSelectedMode('social')}
              className={`flex items-start gap-4 p-4 rounded-xl border-2 cursor-pointer transition-all ${
                selectedMode === 'social'
                  ? 'border-indigo-600 bg-indigo-50/50 shadow-sm'
                  : 'border-slate-200 hover:border-slate-300 bg-white'
              }`}
            >
              <div
                className={`p-2.5 rounded-lg ${
                  selectedMode === 'social' ? 'bg-indigo-600 text-white' : 'bg-slate-100 text-slate-600'
                }`}
              >
                <Users className="w-5 h-5" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <h4 className="font-semibold text-slate-900">Mode Analitik Sosial</h4>
                  {selectedMode === 'social' && (
                    <Badge variant="outline" className="bg-indigo-100 text-indigo-700 border-indigo-200">
                      Terpilih
                    </Badge>
                  )}
                </div>
                <p className="text-xs text-slate-500 mt-1">
                  Menampilkan perbandingan pola pilihan kelompok teman sejawat secara anonim setelah menyelesaikan argumen.
                </p>
              </div>
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-4 border-t border-slate-100">
            <Button variant="outline" onClick={onClose} disabled={isLoading}>
              Batal
            </Button>
            <Button
              className="bg-indigo-600 hover:bg-indigo-700 text-white"
              onClick={handleStart}
              disabled={isLoading}
            >
              {isLoading ? 'Memuat...' : 'Mulai Latihan'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Dialog Konfirmasi Pergantian Mode jika Berbeda dengan Sesi Berjalan */}
      <ConfirmDialog
        isOpen={!!conflictMode}
        onClose={() => {
          if (onCancelConflict) onCancelConflict()
        }}
        onConfirm={() => {
          onSelectMode(selectedMode, true)
        }}
        title="Ubah Mode Belajar?"
        description={`Sesi latihan Anda yang sedang berjalan menggunakan Mode ${conflictMode}. Apakah Anda ingin mengubah mode sesi ini ke Mode ${selectedMode}? Kemajuan argumen Anda tidak akan hilang.`}
        confirmText="Ya, Ubah & Lanjutkan"
        cancelText="Batal"
        isLoading={isLoading}
      />
    </>
  )
}
