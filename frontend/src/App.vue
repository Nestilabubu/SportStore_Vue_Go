<script setup>
import { ref, onMounted, onUnmounted, provide, computed } from "vue";
import { useRouter } from "vue-router";
import Header from "./components/Header.vue";
import Drawer from "./components/CartDrawer.vue";
import {
  getCurrentUser,
  logout as logoutApi,
  refreshSession,
} from "./utils/auth";
import {
  getCart,
  addToCart as apiAddToCart,
  removeFromCart as apiRemoveFromCart,
  updateCartItem as apiUpdateCartItem,
} from "./utils/cart";
import {
  getFavorites,
  addToFavorites as apiAddToFavorites,
  removeFromFavorites as apiRemoveFromFavorites,
} from "./utils/favorites";

const router = useRouter();
const cart = ref([]);
const favorites = ref([]);
const user = ref(null);
const drawerOpen = ref(false);
const isAuthLoading = ref(true);
let refreshInterval = null;

const fetchUser = async () => {
  isAuthLoading.value = true;
  user.value = await getCurrentUser();
  if (user.value) {
    startRefreshTimer(); // запускаем таймер
    cart.value = await getCart();
    favorites.value = await getFavorites();
  } else {
    stopRefreshTimer(); // останавливаем таймер
    cart.value = [];
    favorites.value = [];
  }
  isAuthLoading.value = false;
};

onMounted(() => {
  fetchUser();
});

onUnmounted(() => {
  stopRefreshTimer();
});

const loadInitialData = async () => {
  user.value = await getCurrentUser();
  if (user.value) {
    cart.value = await getCart();
    favorites.value = await getFavorites();
  }
};

onMounted(loadInitialData);

// Методы для работы с корзиной (вызывают API и обновляют локальное состояние)
const addToCart = async (productId, size, quantity = 1) => {
  if (!user.value) {
    alert("Пожалуйста, войдите в систему");
    router.push("/login");
    return;
  }
  await apiAddToCart(productId, size, quantity);
  cart.value = await getCart(); // перезагружаем корзину
};

const startRefreshTimer = () => {
  if (refreshInterval) clearInterval(refreshInterval);
  refreshInterval = setInterval(
    async () => {
      if (user.value) {
        await refreshSession();
      }
    },
    15 * 60 * 1000,
  );
};

const stopRefreshTimer = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval);
    refreshInterval = null;
  }
};

const removeFromCart = async (itemId) => {
  await apiRemoveFromCart(itemId);
  cart.value = await getCart();
};

const updateCartQuantity = async (itemId, quantity) => {
  if (quantity <= 0) {
    await removeFromCart(itemId);
  } else {
    await apiUpdateCartItem(itemId, quantity);
    cart.value = await getCart();
  }
};

// Методы для избранного
const toggleFavorite = async (product) => {
  if (!user.value) {
    alert("Пожалуйста, войдите в систему");
    router.push("/login");
    return;
  }
  const isFav = favorites.value.some((f) => f.id === product.id);
  if (isFav) {
    await apiRemoveFromFavorites(product.id);
  } else {
    await apiAddToFavorites(product.id);
  }
  favorites.value = await getFavorites(); // обновляем список
};

const logout = async () => {
  await logoutApi();
  user.value = null;
  stopRefreshTimer();
  cart.value = [];
  favorites.value = [];
  router.push("/");
};

const refreshCart = async () => {
  if (user.value) {
    cart.value = await getCart();
  } else {
    cart.value = [];
  }
};

const totalPrice = computed(() =>
  cart.value.reduce((acc, item) => acc + item.price * item.quantity, 0),
);
const vatPrice = computed(() => Math.round((totalPrice.value * 5) / 100));

provide("cart", {
  cart,
  totalPrice,
  vatPrice,
  drawerOpen,
  closeDrawer: () => (drawerOpen.value = false),
  openDrawer: () => (drawerOpen.value = true),
  addToCart: async (productId, size, quantity) => {
    await apiAddToCart(productId, size, quantity);
    await refreshCart();
  },
  removeFromCart: async (itemId) => {
    await apiRemoveFromCart(itemId);
    await refreshCart();
  },
  updateCartQuantity: async (itemId, quantity) => {
    if (quantity <= 0) {
      await apiRemoveFromCart(itemId);
    } else {
      await apiUpdateCartItem(itemId, quantity);
    }
    await refreshCart();
  },
  refreshCart, // уже есть
});
provide("favorites", {
  favorites,
  toggleFavorite: async (product) => {
    const isFav = favorites.value.some((f) => f.id === product.id);
    if (isFav) {
      await apiRemoveFromFavorites(product.id);
    } else {
      await apiAddToFavorites(product.id);
    }
    favorites.value = await getFavorites();
  },
  isFavorite: (productId) => favorites.value.some((f) => f.id === productId),
});

provide("auth", {
  user,
  logout,
  isAuthLoading,
  fetchUser,
});
</script>

<template>
  <Drawer v-if="drawerOpen" />
  <div
    class="bg-white w-5/6 m-auto rounded-xl shadow-xl mt-10"
    :class="{ 'opacity-70': drawerOpen }"
  >
    <Header @open-drawer="() => (drawerOpen = true)" />
    <div class="p-10">
      <router-view></router-view>
    </div>
  </div>
</template>
