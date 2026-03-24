import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import { autoAnimatePlugin } from "@formkit/auto-animate/vue";
import "./style.css";
import App from "./App.vue";
import Home from "./pages/Home.vue";
import Favorites from "./pages/Favorites.vue";
import Register from "./pages/Register.vue";
import Login from "./pages/Login.vue";
import Profile from "./pages/Profile.vue";
import Payment from "./pages/Payment.vue";

const app = createApp(App);

const routes = [
  { path: "/", name: "Home", component: Home },
  { path: "/favorites", name: "Favorites", component: Favorites },
  { path: "/register", name: "Register", component: Register },
  { path: "/login", name: "Login", component: Login },
  { path: "/profile", name: "Profile", component: Profile },
  { path: '/payment', name: 'Payment', component: Payment },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

app.use(router);
app.use(autoAnimatePlugin);
app.mount("#app");
