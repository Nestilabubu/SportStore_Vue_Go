import axios from "axios";

const getBaseURL = () => {
  if (import.meta.env.PROD) {
    return (
      import.meta.env.VITE_API_URL ||
      "https://sportstorevuego-production.up.railway.app/api"
    );
  }
  return "http://localhost:8080/api";
};

const api = axios.create({
  baseURL: getBaseURL(),
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

// Интерцептор для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      console.log("Unauthorized, please login");
    }
    return Promise.reject(error);
  },
);

export default api;
