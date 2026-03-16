import api from "./api";

export async function getStatistics() {
  const response = await api.get("/statistics");
  return response.data;
}
