import api from "./api";

export async function register(userData) {
  const response = await api.post("/register", userData);
  return response.data;
}

export async function login(credentials) {
  const response = await api.post("/login", credentials);
  return response.data;
}

export async function logout() {
  await api.post("/logout");
}

export async function getCurrentUser() {
  try {
    const response = await api.get("/profile", {
      validateStatus: (status) => status < 500,
    });
    if (response.status === 401) return null;
    return response.data;
  } catch {
    return null;
  }
}

export async function updateProfile(data) {
  const response = await api.put("/profile", data);
  return response.data;
}
