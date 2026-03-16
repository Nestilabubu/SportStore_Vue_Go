import api from "./api";

export async function getCart() {
  const response = await api.get("/cart");
  return response.data;
}

export async function addToCart(productId, size, quantity = 1) {
  const response = await api.post("/cart", { productId, size, quantity });
  return response.data;
}

export async function updateCartItem(itemId, quantity) {
  const response = await api.put(`/cart/${itemId}`, { quantity });
  return response.data;
}

export async function removeFromCart(itemId) {
  await api.delete(`/cart/${itemId}`);
}
