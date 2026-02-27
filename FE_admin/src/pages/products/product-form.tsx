import { useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { ArrowLeft, Loader2 } from 'lucide-react'
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
import { productsApi } from '@/services/products'
import type { ProductCreatePayload, ProductStatus } from '@/types'

interface ProductFormValues {
  sku: string
  name: string
  description: string
  price: string
  stock: number
  status: ProductStatus
  image_url: string
}

export default function ProductFormPage() {
  const { id } = useParams<{ id: string }>()
  const isEdit = !!id
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: product, isLoading: loadingProduct } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productsApi.get(id!),
    enabled: isEdit,
  })

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
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
        image_url: product.image_url ?? '',
      })
    }
  }, [product, reset])

  const createMutation = useMutation({
    mutationFn: (payload: ProductCreatePayload) => productsApi.create(payload),
    onSuccess: () => {
      toast.success('Tạo sản phẩm thành công')
      queryClient.invalidateQueries({ queryKey: ['products'] })
      navigate('/products')
    },
    onError: () => toast.error('Tạo sản phẩm thất bại'),
  })

  const updateMutation = useMutation({
    mutationFn: (payload: ProductCreatePayload) => productsApi.update(id!, payload),
    onSuccess: () => {
      toast.success('Cập nhật sản phẩm thành công')
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['product', id] })
      navigate('/products')
    },
    onError: () => toast.error('Cập nhật sản phẩm thất bại'),
  })

  const onSubmit = (data: ProductFormValues) => {
    const payload: ProductCreatePayload = {
      sku: data.sku,
      name: data.name,
      description: data.description || undefined,
      price: parseFloat(data.price).toFixed(2),
      stock: data.stock,
      status: data.status,
      image_url: data.image_url || undefined,
    }
    if (isEdit) {
      updateMutation.mutate(payload)
    } else {
      createMutation.mutate(payload)
    }
  }

  const isSubmitting = createMutation.isPending || updateMutation.isPending

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
                  <Select
                    value={watch('status')}
                    onValueChange={(v: string) => setValue('status', v as ProductStatus)}
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
                <Label htmlFor="image_url">URL hình ảnh</Label>
                <Input
                  id="image_url"
                  placeholder="https://..."
                  {...register('image_url')}
                />
              </div>
              {watch('image_url') && (
                <img
                  src={watch('image_url')}
                  alt="Preview"
                  className="h-32 w-32 rounded-md object-cover border"
                  onError={(e) => { (e.target as HTMLImageElement).style.display = 'none' }}
                />
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
                  Đang lưu...
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
