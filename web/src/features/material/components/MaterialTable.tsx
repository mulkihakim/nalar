import { useState } from 'react'
import { Plus, Search, Edit3, Trash2, BookOpen, Layers, ChevronDown, ChevronUp } from 'lucide-react'
import type { Material, Argument } from '../types'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { ConfirmDialog } from '../../../components/ConfirmDialog'
import { MaterialModal } from './MaterialModal'
import { ArgumentModal } from './ArgumentModal'
import { useDeleteMaterial, useDeleteArgument } from '../api'

interface MaterialTableProps {
  materials: Material[]
  isLoading: boolean
  currentUserRole: string
  currentUserId: number
}

export function MaterialTable({
  materials,
  isLoading,
  currentUserRole,
  currentUserId,
}: MaterialTableProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [isMaterialModalOpen, setIsMaterialModalOpen] = useState(false)
  const [selectedMaterial, setSelectedMaterial] = useState<Material | null>(null)

  const [isArgumentModalOpen, setIsArgumentModalOpen] = useState(false)
  const [activeMaterialId, setActiveMaterialId] = useState<number>(0)
  const [selectedArgument, setSelectedArgument] = useState<Argument | null>(null)

  const [expandedMaterialId, setExpandedMaterialId] = useState<number | null>(null)
  const [expandedContentIds, setExpandedContentIds] = useState<Set<number>>(new Set())

  const toggleContentExpand = (id: number) => {
    setExpandedContentIds((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  // Confirm delete states
  const [materialToDelete, setMaterialToDelete] = useState<Material | null>(null)
  const [argumentToDelete, setArgumentToDelete] = useState<{
    materialId: number
    argument: Argument
  } | null>(null)
  const [actionError, setActionError] = useState<string | null>(null)

  const deleteMaterialMutation = useDeleteMaterial()
  const deleteArgumentMutation = useDeleteArgument()

  const filteredMaterials = materials.filter((m) =>
    m.title.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const handleCreateMaterial = () => {
    setSelectedMaterial(null)
    setIsMaterialModalOpen(true)
  }

  const handleEditMaterial = (m: Material) => {
    setSelectedMaterial(m)
    setIsMaterialModalOpen(true)
  }

  const handleDeleteMaterialConfirm = async () => {
    if (!materialToDelete) return
    setActionError(null)
    try {
      await deleteMaterialMutation.mutateAsync(materialToDelete.id)
      setMaterialToDelete(null)
    } catch (err: any) {
      setActionError(err.message || 'Gagal menghapus materi')
    }
  }

  const handleAddArgument = (materialId: number) => {
    setActiveMaterialId(materialId)
    setSelectedArgument(null)
    setIsArgumentModalOpen(true)
  }

  const handleEditArgument = (materialId: number, arg: Argument) => {
    setActiveMaterialId(materialId)
    setSelectedArgument(arg)
    setIsArgumentModalOpen(true)
  }

  const handleDeleteArgumentConfirm = async () => {
    if (!argumentToDelete) return
    setActionError(null)
    try {
      await deleteArgumentMutation.mutateAsync({
        materialId: argumentToDelete.materialId,
        argumentId: argumentToDelete.argument.id,
      })
      setArgumentToDelete(null)
    } catch (err: any) {
      setActionError(err.message || 'Gagal menghapus argumen')
    }
  }

  const canManageMaterial = (m: Material) => {
    if (currentUserRole === 'admin') return true
    if (currentUserRole === 'asesor' && m.owner_id === currentUserId) return true
    return false
  }

  return (
    <div className="space-y-4">
      {actionError && (
        <div className="p-3 text-xs bg-red-50 border border-red-200 text-red-700 rounded-lg flex items-center justify-between">
          <span>{actionError}</span>
          <button
            type="button"
            onClick={() => setActionError(null)}
            className="text-red-500 hover:text-red-800 font-bold ml-2 cursor-pointer"
          >
            &times;
          </button>
        </div>
      )}

      {/* Action Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3 items-stretch sm:items-center justify-between">
        <div className="relative flex-1 max-w-md">
          <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <Input
            placeholder="Cari judul materi..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-9 bg-white"
          />
        </div>

        <Button onClick={handleCreateMaterial} className="gap-1.5 cursor-pointer">
          <Plus className="w-4 h-4" />
          <span>Buat Materi Baru</span>
        </Button>
      </div>

      {/* Materials List */}
      <div className="space-y-3">
        {isLoading ? (
          <div className="bg-white p-8 rounded-xl border border-slate-200 text-center text-slate-400">
            Memuat daftar materi...
          </div>
        ) : filteredMaterials.length === 0 ? (
          <div className="bg-white p-8 rounded-xl border border-slate-200 text-center text-slate-400">
            Belum ada materi. Klik "Buat Materi Baru" untuk memulai.
          </div>
        ) : (
          filteredMaterials.map((m) => {
            const isOwnerOrAdmin = canManageMaterial(m)
            const isExpanded = expandedMaterialId === m.id
            const args = m.arguments || []

            return (
              <div
                key={m.id}
                className="bg-white rounded-xl border border-slate-200 shadow-xs overflow-hidden transition-all"
              >
                {/* Header Card */}
                <div className="p-4 sm:p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-slate-100">
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <BookOpen className="w-4 h-4 text-indigo-600" />
                      <h3 className="font-bold text-slate-900 text-base">{m.title}</h3>
                    </div>
                    <div className="text-xs text-slate-500 flex items-center gap-3">
                      <span>Pemilik: {m.owner?.name || `User #${m.owner_id}`}</span>
                      <span>•</span>
                      <span className="font-medium text-indigo-700 bg-indigo-50 px-2 py-0.5 rounded-full">
                        {args.length} Argumen
                      </span>
                    </div>
                  </div>

                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="xs"
                      onClick={() => setExpandedMaterialId(isExpanded ? null : m.id)}
                      className="gap-1 cursor-pointer text-slate-700"
                    >
                      <Layers className="w-3 h-3 text-slate-500" />
                      <span>{isExpanded ? 'Tutup Argumen' : 'Lihat Argumen'}</span>
                    </Button>

                    {isOwnerOrAdmin && (
                      <>
                        <Button
                          variant="outline"
                          size="xs"
                          onClick={() => handleAddArgument(m.id)}
                          className="gap-1 cursor-pointer text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50"
                        >
                          <Plus className="w-3 h-3" />
                          <span>Argumen</span>
                        </Button>
                        <Button
                          variant="outline"
                          size="xs"
                          onClick={() => handleEditMaterial(m)}
                          className="cursor-pointer"
                        >
                          <Edit3 className="w-3 h-3" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="xs"
                          onClick={() => setMaterialToDelete(m)}
                          className="text-red-500 hover:text-red-700 hover:bg-red-50 cursor-pointer"
                          title="Hapus Materi"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                        </Button>
                      </>
                    )}
                  </div>
                </div>

                {/* Content Preview */}
                <div className="px-5 py-3.5 bg-slate-50/70 border-b border-slate-100 text-xs sm:text-sm text-slate-700 leading-relaxed">
                  {(() => {
                    const isContentExpanded = expandedContentIds.has(m.id)
                    const isLongContent = m.content.length > 450
                    const displayText =
                      isLongContent && !isContentExpanded
                        ? `${m.content.slice(0, 420).trim()}...`
                        : m.content

                    return (
                      <div className="space-y-2">
                        <p className="whitespace-pre-line break-words">
                          {displayText}
                        </p>
                        {isLongContent && (
                          <div>
                            <button
                              type="button"
                              onClick={() => toggleContentExpand(m.id)}
                              className="inline-flex items-center gap-1 text-xs font-semibold text-indigo-600 hover:text-indigo-800 transition-colors cursor-pointer"
                            >
                              {isContentExpanded ? (
                                <>
                                  <ChevronUp className="w-3.5 h-3.5" />
                                  <span>Tampilkan Lebih Sedikit</span>
                                </>
                              ) : (
                                <>
                                  <ChevronDown className="w-3.5 h-3.5" />
                                  <span>Baca Selengkapnya</span>
                                </>
                              )}
                            </button>
                          </div>
                        )}
                      </div>
                    )
                  })()}
                </div>

                {/* Expandable Arguments Panel */}
                {isExpanded && (
                  <div className="p-4 sm:p-5 bg-slate-50/80 space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-bold uppercase tracking-wider text-slate-500">
                        Daftar Argumen Materi Ini ({args.length})
                      </span>
                      {isOwnerOrAdmin && (
                        <Button
                          size="xs"
                          onClick={() => handleAddArgument(m.id)}
                          className="gap-1 cursor-pointer"
                        >
                          <Plus className="w-3 h-3" />
                          <span>Tambah Argumen</span>
                        </Button>
                      )}
                    </div>

                    {args.length === 0 ? (
                      <div className="p-4 text-center text-xs text-slate-400 bg-white rounded-lg border border-slate-200">
                        Belum ada argumen untuk materi ini.
                      </div>
                    ) : (
                      <div className="space-y-2">
                        {args.map((a, idx) => (
                          <div
                            key={a.id}
                            className="bg-white p-3.5 rounded-lg border border-slate-200 shadow-xs space-y-2"
                          >
                            <div className="flex items-start justify-between gap-2">
                              <div>
                                <span className="text-xs font-bold text-indigo-600 mr-2">
                                  #{idx + 1}
                                </span>
                                <span className="text-sm font-semibold text-slate-900">
                                  {a.claim_text}
                                </span>
                              </div>
                              {isOwnerOrAdmin && (
                                <div className="flex items-center gap-1 shrink-0">
                                  <Button
                                    variant="outline"
                                    size="xs"
                                    onClick={() => handleEditArgument(m.id, a)}
                                    className="cursor-pointer"
                                  >
                                    <Edit3 className="w-3 h-3" />
                                  </Button>
                                  <Button
                                    variant="ghost"
                                    size="xs"
                                    onClick={() =>
                                      setArgumentToDelete({ materialId: m.id, argument: a })
                                    }
                                    className="text-red-500 hover:text-red-700 hover:bg-red-50 cursor-pointer"
                                    title="Hapus Argumen"
                                  >
                                    <Trash2 className="w-3.5 h-3.5" />
                                  </Button>
                                </div>
                              )}
                            </div>

                            {/* Options Preview */}
                            {a.options && (
                              <div className="grid grid-cols-1 md:grid-cols-2 gap-2 pt-1 border-t border-slate-100 text-[11px]">
                                <div className="space-y-1">
                                  <span className="font-semibold text-amber-800">Ground:</span>
                                  {a.options
                                    .filter((o) => o.type === 'ground')
                                    .map((o) => (
                                      <div
                                        key={o.id}
                                        className={`px-2 py-0.5 rounded ${
                                          o.is_correct
                                            ? 'bg-emerald-50 text-emerald-800 font-semibold border border-emerald-200'
                                            : 'text-slate-600'
                                        }`}
                                      >
                                        {o.is_correct ? '✓ ' : '• '}
                                        {o.text}
                                      </div>
                                    ))}
                                </div>
                                <div className="space-y-1">
                                  <span className="font-semibold text-indigo-800">Warrant:</span>
                                  {a.options
                                    .filter((o) => o.type === 'warrant')
                                    .map((o) => (
                                      <div
                                        key={o.id}
                                        className={`px-2 py-0.5 rounded ${
                                          o.is_correct
                                            ? 'bg-emerald-50 text-emerald-800 font-semibold border border-emerald-200'
                                            : 'text-slate-600'
                                        }`}
                                      >
                                        {o.is_correct ? '✓ ' : '• '}
                                        {o.text}
                                      </div>
                                    ))}
                                </div>
                              </div>
                            )}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                )}
              </div>
            )
          })
        )}
      </div>

      <MaterialModal
        isOpen={isMaterialModalOpen}
        onClose={() => setIsMaterialModalOpen(false)}
        materialToEdit={selectedMaterial}
      />

      <ArgumentModal
        isOpen={isArgumentModalOpen}
        onClose={() => setIsArgumentModalOpen(false)}
        materialId={activeMaterialId}
        argumentToEdit={selectedArgument}
      />

      {/* Confirm Dialog for Material Deletion */}
      <ConfirmDialog
        isOpen={!!materialToDelete}
        onClose={() => setMaterialToDelete(null)}
        onConfirm={handleDeleteMaterialConfirm}
        title="Hapus Materi"
        description={`Apakah Anda yakin ingin menghapus materi "${materialToDelete?.title}" beserta seluruh argumen di dalamnya?`}
        confirmText="Hapus Materi"
        variant="destructive"
        isLoading={deleteMaterialMutation.isPending}
      />

      {/* Confirm Dialog for Argument Deletion */}
      <ConfirmDialog
        isOpen={!!argumentToDelete}
        onClose={() => setArgumentToDelete(null)}
        onConfirm={handleDeleteArgumentConfirm}
        title="Hapus Argumen"
        description={`Apakah Anda yakin ingin menghapus argumen "${argumentToDelete?.argument.claim_text}"?`}
        confirmText="Hapus Argumen"
        variant="destructive"
        isLoading={deleteArgumentMutation.isPending}
      />
    </div>
  )
}
