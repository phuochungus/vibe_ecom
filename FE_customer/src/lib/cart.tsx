import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react'
import type { Product, CartItem } from '@/types'

interface CartContextValue {
  items: CartItem[]
  itemCount: number
  totalAmount: number
  addItem: (product: Product, quantity?: number) => void
  updateQuantity: (productId: string, quantity: number) => void
  removeItem: (productId: string) => void
  clearCart: () => void
  isInCart: (productId: string) => boolean
  getItem: (productId: string) => CartItem | undefined
}

const CartContext = createContext<CartContextValue | null>(null)

const STORAGE_KEY = 'shopping_cart'

function getCartFromStorage(): CartItem[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : []
  } catch {
    return []
  }
}

function saveCartToStorage(items: CartItem[]) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(items))
  } catch (error) {
    console.error('Failed to save cart to localStorage:', error)
  }
}

export function CartProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<CartItem[]>(getCartFromStorage)

  // Persist to localStorage whenever items change
  useEffect(() => {
    saveCartToStorage(items)
  }, [items])

  const addItem = useCallback((product: Product, quantity: number = 1) => {
    // BR-PROD-03: Filter out DISCONTINUED products
    if (product.status === 'DISCONTINUED') {
      throw new Error('Sản phẩm này đã ngừng bán')
    }

    // Check stock
    setItems((prev) => {
      const existing = prev.find((item) => item.product_id === product.id)
      if (existing) {
        const newQuantity = existing.quantity + quantity
        if (newQuantity > product.stock) {
          throw new Error(`Chỉ còn ${product.stock} sản phẩm trong kho`)
        }
        return prev.map((item) =>
          item.product_id === product.id
            ? { ...item, quantity: newQuantity }
            : item
        )
      } else {
        if (quantity > product.stock) {
          throw new Error(`Chỉ còn ${product.stock} sản phẩm trong kho`)
        }
        return [...prev, { product_id: product.id, product, quantity }]
      }
    })
  }, [])

  const updateQuantity = useCallback((productId: string, quantity: number) => {
    if (quantity <= 0) {
      setItems((prev) => prev.filter((item) => item.product_id !== productId))
      return
    }

    setItems((prev) =>
      prev.map((item) => {
        if (item.product_id === productId) {
          if (quantity > item.product.stock) {
            throw new Error(`Chỉ còn ${item.product.stock} sản phẩm trong kho`)
          }
          return { ...item, quantity }
        }
        return item
      })
    )
  }, [])

  const removeItem = useCallback((productId: string) => {
    setItems((prev) => prev.filter((item) => item.product_id !== productId))
  }, [])

  const clearCart = useCallback(() => {
    setItems([])
  }, [])

  const isInCart = useCallback((productId: string) => {
    return items.some((item) => item.product_id === productId)
  }, [items])

  const getItem = useCallback((productId: string) => {
    return items.find((item) => item.product_id === productId)
  }, [items])

  const itemCount = items.reduce((sum, item) => sum + item.quantity, 0)
  const totalAmount = items.reduce((sum, item) => {
    const price = parseFloat(item.product.price)
    return sum + price * item.quantity
  }, 0)

  return (
    <CartContext.Provider
      value={{
        items,
        itemCount,
        totalAmount,
        addItem,
        updateQuantity,
        removeItem,
        clearCart,
        isInCart,
        getItem,
      }}
    >
      {children}
    </CartContext.Provider>
  )
}

export function useCart() {
  const ctx = useContext(CartContext)
  if (!ctx) throw new Error('useCart must be used within CartProvider')
  return ctx
}
