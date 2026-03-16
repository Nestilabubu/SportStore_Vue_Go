import api from "./api";

export async function getFavorites() {
  const response = await api.get("/favorites");
  return response.data;
}

export async function addToFavorites(productId) {
  const response = await api.post("/favorites", { productId });
  return response.data;
}

export async function removeFromFavorites(productId) {
  await api.delete(`/favorites/${productId}`);
}

export async function isItemInFavorites(productId) {
  const favs = await getFavorites();
  return favs.some((f) => f.id === productId);
}
