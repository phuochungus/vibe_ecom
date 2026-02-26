import { createContext, useContext, useMemo, useState } from "react";

export type CartItem = {
  productId: string;
  name: string;
  price: string;
  quantity: number;
};

type CartContextValue = {
  items: CartItem[];
  addItem: (item: Omit<CartItem, "quantity">) => void;
  removeItem: (productId: string) => void;
  clear: () => void;
  totalAmount: string;
};

const CartContext = createContext<CartContextValue | undefined>(undefined);

const toAmount = (value: number): string => value.toFixed(2);

export const CartProvider = ({ children }: { children: React.ReactNode }) => {
  const [items, setItems] = useState<CartItem[]>([]);

  const addItem = (item: Omit<CartItem, "quantity">) => {
    setItems((prev) => {
      const idx = prev.findIndex((it) => it.productId === item.productId);
      if (idx === -1) {
        return [...prev, { ...item, quantity: 1 }];
      }

      const next = [...prev];
      next[idx] = { ...next[idx], quantity: next[idx].quantity + 1 };
      return next;
    });
  };

  const removeItem = (productId: string) => {
    setItems((prev) => prev.filter((it) => it.productId !== productId));
  };

  const clear = () => setItems([]);

  const totalAmount = useMemo(() => {
    const total = items.reduce((acc, item) => acc + Number(item.price) * item.quantity, 0);
    return toAmount(total);
  }, [items]);

  const value = useMemo(
    () => ({ items, addItem, removeItem, clear, totalAmount }),
    [items, totalAmount],
  );

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
};

export const useCart = (): CartContextValue => {
  const ctx = useContext(CartContext);
  if (!ctx) {
    throw new Error("useCart must be used within CartProvider");
  }
  return ctx;
};
