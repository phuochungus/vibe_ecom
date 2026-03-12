import { useEffect, useState, type ChangeEvent } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Controller, useForm } from 'react-hook-form'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { ArrowLeft, ImagePlus, Loader2, Upload, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import { getErrorMessage } from '@/lib/api'
import { productsApi } from '@/services/products'
import type { ProductCreatePayload, ProductStatus } from '@/types'

interface ProductFormValues {
  sku: string
  name: string
  description: string
  price: string
  stock: number
  status: ProductStatus
}

export default function ProductFormPage() {
  const { id } = useParams<{ id: string }>()
  const isEdit = !!id
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [selectedImageFile, setSelectedImageFile] = useState<File | null>(null)
  const [localPreviewUrl, setLocalPreviewUrl] = useState<string | null>(null)
  const [uploadedImageUrl, setUploadedImageUrl] = useState<string | null>(null)

  const { data: product, isLoading: loadingProduct } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productsApi.get(id!),
    enabled: isEdit,
  })

  const {
    control,
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ProductFormValues>({
    defaultValues: {
      status: 'ACTIVE',
      stock: 0,
    },
  })

  useEffect(() => {
    if (product) {
      reset({
        sku: product.sku,
        name: product.name,
        description: product.description ?? '',
        price: product.price,
        stock: product.stock,
        status: product.status as ProductStatus,
      })
    }
  }, [product, reset])

  useEffect(() => {
    return () => {
      if (localPreviewUrl) {
        URL.revokeObjectURL(localPreviewUrl)
      }
    }
  }, [localPreviewUrl])

  const createMutation = useMutation({
    mutationFn: (payload: ProductCreatePayload) => productsApi.create(payload),
    onSuccess: () => {
      toast.success('Tạo sản phẩm thành công')
      queryClient.invalidateQueries({ queryKey: ['products'] })
      navigate('/products')
    },
    onError: (error) => toast.error(getErrorMessage(error) || 'Tạo sản phẩm thất bại'),
  })

  const updateMutation = useMutation({
    mutationFn: (payload: ProductCreatePayload) => productsApi.update(id!, payload),
    onSuccess: () => {
      toast.success('Cập nhật sản phẩm thành công')
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['product', id] })
      navigate('/products')
    },
    onError: (error) => toast.error(getErrorMessage(error) || 'Cập nhật sản phẩm thất bại'),
  })

  const uploadImageMutation = useMutation({
    mutationFn: (file: File) => productsApi.uploadImage(file),
    onError: (error) => toast.error(getErrorMessage(error) || 'Tải ảnh lên thất bại'),
  })

  const handleImageFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) {
      return
    }

    if (localPreviewUrl) {
      URL.revokeObjectURL(localPreviewUrl)
    }

    const previewUrl = URL.createObjectURL(file)
    setSelectedImageFile(file)
    setLocalPreviewUrl(previewUrl)
  }

  const clearSelectedImage = () => {
    if (localPreviewUrl) {
      URL.revokeObjectURL(localPreviewUrl)
    }
    setSelectedImageFile(null)
    setLocalPreviewUrl(null)
  }

  const onSubmit = async (data: ProductFormValues) => {
    try {
      let imageURL = uploadedImageUrl || product?.image_url || undefined

      if (selectedImageFile) {
        const upload = await uploadImageMutation.mutateAsync(selectedImageFile)
        imageURL = upload.url
        setUploadedImageUrl(upload.url)
      }

      const payload: ProductCreatePayload = {
        sku: data.sku,
        name: data.name,
        description: data.description || undefined,
        price: Number.parseInt(data.price, 10).toString(),
        stock: data.stock,
        status: data.status,
        image_url: imageURL,
      }

      if (isEdit) {
        await updateMutation.mutateAsync(payload)
        return
      }

      await createMutation.mutateAsync(payload)
    } catch {
      return
    }
  }

  const isSubmitting =
    createMutation.isPending || updateMutation.isPending || uploadImageMutation.isPending
  const imagePreview = localPreviewUrl || uploadedImageUrl || product?.image_url || ''

  return (
    <div className="max-w-2xl space-y-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={() => navigate('/products')}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-bold">{isEdit ? 'Sửa sản phẩm' : 'Thêm sản phẩm mới'}</h1>
          <p className="text-sm text-muted-foreground">
            {isEdit ? `ID: ${id}` : 'Điền thông tin sản phẩm mới'}
          </p>
        </div>
      </div>

      {isEdit && loadingProduct ? (
        <div className="space-y-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-10 w-full" />
          ))}
        </div>
      ) : (
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Thông tin cơ bản</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="sku">
                    SKU <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="sku"
                    placeholder="CLUB-001"
                    {...register('sku', { required: 'SKU là bắt buộc' })}
                    disabled={isEdit}
                  />
                  {errors.sku && <p className="text-xs text-destructive">{errors.sku.message}</p>}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="status">Trạng thái</Label>
                  <Controller
                    name="status"
                    control={control}
                    render={({ field }) => (
                      <Select
                        value={field.value}
                        onValueChange={(value: string) => field.onChange(value as ProductStatus)}
                      >
                        <SelectTrigger id="status">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="ACTIVE">Đang bán</SelectItem>
                          <SelectItem value="INACTIVE">Tạm dừng</SelectItem>
                          <SelectItem value="DISCONTINUED">Ngừng bán</SelectItem>
                        </SelectContent>
                      </Select>
                    )}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="name">
                  Tên sản phẩm <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="name"
                  placeholder="Driver X Pro"
                  {...register('name', { required: 'Tên sản phẩm là bắt buộc' })}
                />
                {errors.name && <p className="text-xs text-destructive">{errors.name.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">Mô tả</Label>
                <Textarea
                  id="description"
                  placeholder="Mô tả chi tiết về sản phẩm..."
                  rows={3}
                  {...register('description')}
                />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-base">Giá & Tồn kho</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="price">
                    Giá bán (VND) <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="price"
                    type="number"
                    placeholder="1590000"
                    {...register('price', {
                      required: 'Giá bán là bắt buộc',
                      min: { value: 1, message: 'Giá phải lớn hơn 0 (BR-PROD-01)' },
                    })}
                  />
                  {errors.price && (
                    <p className="text-xs text-destructive">{errors.price.message}</p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="stock">
                    Tồn kho <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="stock"
                    type="number"
                    placeholder="100"
                    min="0"
                    {...register('stock', {
                      valueAsNumber: true,
                      required: 'Tồn kho là bắt buộc',
                      min: { value: 0, message: 'Tồn kho không được âm (BR-PROD-02)' },
                    })}
                  />
                  {errors.stock && (
                    <p className="text-xs text-destructive">{errors.stock.message}</p>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-base">Hình ảnh</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="image_file">Tải ảnh sản phẩm</Label>
                <Input
                  id="image_file"
                  type="file"
                  accept="image/*"
                  onChange={handleImageFileChange}
                  disabled={isSubmitting}
                />
                <p className="text-xs text-muted-foreground">
                  Hỗ trợ JPG, PNG, WebP, GIF, AVIF. Kích thước tối đa 10MB.
                </p>
              </div>

              {selectedImageFile ? (
                <div className="flex items-center justify-between rounded-md border border-dashed px-3 py-2 text-sm">
                  <div className="flex items-center gap-2 text-muted-foreground">
                    <Upload className="h-4 w-4" />
                    <span className="truncate">{selectedImageFile.name}</span>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={clearSelectedImage}
                    disabled={isSubmitting}
                  >
                    <X className="mr-1 h-4 w-4" />
                    Bỏ chọn
                  </Button>
                </div>
              ) : null}

              {imagePreview ? (
                <div className="space-y-2">
                  <p className="text-sm font-medium">
                    {selectedImageFile ? 'Ảnh xem trước' : 'Ảnh hiện tại'}
                  </p>
                  <img
                    src={imagePreview}
                    alt="Preview"
                    className="h-40 w-40 rounded-md border object-cover"
                    onError={(e) => {
                      const image = e.target as HTMLImageElement
                      image.style.display = 'none'
                    }}
                  />
                </div>
              ) : (
                <div className="flex h-40 w-40 flex-col items-center justify-center rounded-md border border-dashed text-center text-sm text-muted-foreground">
                  <ImagePlus className="mb-2 h-5 w-5" />
                  Chưa có ảnh sản phẩm
                </div>
              )}
            </CardContent>
          </Card>

          <Separator />

          <div className="flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={() => navigate('/products')}>
              Hủy
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  {uploadImageMutation.isPending ? 'Đang tải ảnh...' : 'Đang lưu...'}
                </>
              ) : isEdit ? (
                'Cập nhật'
              ) : (
                'Tạo sản phẩm'
              )}
            </Button>
          </div>
        </form>
      )}
    </div>
  )
}
