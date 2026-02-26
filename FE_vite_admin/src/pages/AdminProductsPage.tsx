import { useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";

import { apiClient } from "../api";
import type { Product, ProductStatus } from "../api";
import { getEnvelopeItems } from "../utils/api";

type ProductFormState = {
  sku: string;
  name: string;
  description: string;
  price: string;
  stock: number;
  status: ProductStatus;
  image_url: string;
};

const emptyForm: ProductFormState = {
  sku: "",
  name: "",
  description: "",
  price: "",
  stock: 0,
  status: "ACTIVE",
  image_url: "",
};

export const AdminProductsPage = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [statusFilter, setStatusFilter] = useState<ProductStatus | "ALL">("ALL");
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [form, setForm] = useState<ProductFormState>(emptyForm);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchProducts = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.products.productsGet(
        undefined,
        statusFilter === "ALL" ? undefined : statusFilter,
        undefined,
        undefined,
        1,
        100,
      );
      const items = getEnvelopeItems<Product>(response.data);
      setProducts(items);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Không tải được danh sách sản phẩm.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchProducts();
  }, [statusFilter]);

  const onSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setError(null);

    try {
      if (editingProduct) {
        await apiClient.adminProducts.adminProductsProductIdPatch(editingProduct.id, {
          name: form.name,
          description: form.description,
          price: form.price,
          stock: Number(form.stock),
          status: form.status,
          image_url: form.image_url || null,
        });
      } else {
        await apiClient.adminProducts.adminProductsPost({
          sku: form.sku,
          name: form.name,
          description: form.description,
          price: form.price,
          stock: Number(form.stock),
          status: form.status,
          image_url: form.image_url || null,
        });
      }
      setEditingProduct(null);
      setForm(emptyForm);
      await fetchProducts();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Không thể lưu sản phẩm.");
    }
  };

  const onEdit = (product: Product) => {
    setEditingProduct(product);
    setForm({
      sku: product.sku,
      name: product.name,
      description: product.description ?? "",
      price: product.price,
      stock: product.stock,
      status: product.status,
      image_url: product.image_url ?? "",
    });
  };

  const onDelete = async (productId: string) => {
    setError(null);
    try {
      await apiClient.adminProducts.adminProductsProductIdDelete(productId);
      await fetchProducts();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Không thể xóa sản phẩm.");
    }
  };

  const title = useMemo(() => (editingProduct ? "Cập nhật sản phẩm" : "Tạo sản phẩm mới"), [editingProduct]);

  return (
    <section className="page-grid">
      <article className="card">
        <div className="section-header">
          <h2>Quản lý sản phẩm</h2>
          <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value as ProductStatus | "ALL")}>
            <option value="ALL">Tất cả</option>
            <option value="ACTIVE">ACTIVE</option>
            <option value="INACTIVE">INACTIVE</option>
            <option value="DISCONTINUED">DISCONTINUED</option>
          </select>
        </div>
        {loading && <p className="muted">Đang tải dữ liệu...</p>}
        {error && <p className="error">{error}</p>}
        <table>
          <thead>
            <tr>
              <th>SKU</th>
              <th>Tên</th>
              <th>Giá</th>
              <th>Tồn kho</th>
              <th>Trạng thái</th>
              <th>Thao tác</th>
            </tr>
          </thead>
          <tbody>
            {products.map((product) => (
              <tr key={product.id}>
                <td>{product.sku}</td>
                <td>{product.name}</td>
                <td>{product.price}</td>
                <td>{product.stock}</td>
                <td>{product.status}</td>
                <td className="actions-cell">
                  <button onClick={() => onEdit(product)}>Sửa</button>
                  <button className="danger" onClick={() => void onDelete(product.id)}>
                    Xóa
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </article>

      <article className="card">
        <h3>{title}</h3>
        <form className="form-grid" onSubmit={onSubmit}>
          {!editingProduct && (
            <label>
              SKU
              <input value={form.sku} onChange={(event) => setForm((prev) => ({ ...prev, sku: event.target.value }))} required />
            </label>
          )}
          <label>
            Tên
            <input value={form.name} onChange={(event) => setForm((prev) => ({ ...prev, name: event.target.value }))} required />
          </label>
          <label>
            Mô tả
            <textarea value={form.description} onChange={(event) => setForm((prev) => ({ ...prev, description: event.target.value }))} rows={3} />
          </label>
          <label>
            Giá
            <input value={form.price} onChange={(event) => setForm((prev) => ({ ...prev, price: event.target.value }))} required />
          </label>
          <label>
            Tồn kho
            <input
              type="number"
              min={0}
              value={form.stock}
              onChange={(event) => setForm((prev) => ({ ...prev, stock: Number(event.target.value) }))}
              required
            />
          </label>
          <label>
            Trạng thái
            <select value={form.status} onChange={(event) => setForm((prev) => ({ ...prev, status: event.target.value as ProductStatus }))}>
              <option value="ACTIVE">ACTIVE</option>
              <option value="INACTIVE">INACTIVE</option>
              <option value="DISCONTINUED">DISCONTINUED</option>
            </select>
          </label>
          <label>
            Image URL
            <input value={form.image_url} onChange={(event) => setForm((prev) => ({ ...prev, image_url: event.target.value }))} />
          </label>
          <div className="actions-row">
            <button type="submit">{editingProduct ? "Lưu thay đổi" : "Tạo sản phẩm"}</button>
            {editingProduct && (
              <button
                type="button"
                onClick={() => {
                  setEditingProduct(null);
                  setForm(emptyForm);
                }}
              >
                Hủy
              </button>
            )}
          </div>
        </form>
      </article>
    </section>
  );
};
