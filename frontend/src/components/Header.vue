<script setup>
import { inject, computed } from "vue";
import { useRouter } from "vue-router";

const { totalPrice, openDrawer } = inject("cart");
const { favorites } = inject("favorites");
const { user, logout } = inject("auth");
const router = useRouter();

const favoritesCount = computed(() => favorites.value?.length || 0);

const handleLogout = async () => {
  await logout();
  router.push("/");
};
</script>

<template>
  <header class="flex justify-between border-b border-slate-200 px-10 py-8">
    <router-link to="/">
      <div class="flex items-center cursor-pointer gap-4">
        <img src="/logo.png" alt="logo" class="w-20" />
        <div>
          <h2 class="text-xl font-bold uppercase">SportStyle</h2>
          <p class="text-slate-400">Магазин спортивных костюмов</p>
        </div>
      </div>
    </router-link>

    <ul class="flex items-center gap-10">
      <template v-if="user">
        <li>
          <router-link to="/profile">
            <div
              class="flex items-center gap-3 text-gray-600 hover:text-black transition-colors group"
            >
              <div
                class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center group-hover:bg-blue-200 transition"
              >
                <span class="text-blue-600 font-medium">{{
                  user.fullName?.[0] || "U"
                }}</span>
              </div>
              <div>
                <p class="font-medium">
                  {{ user.fullName?.split(" ")[0] || "Профиль" }}
                </p>
                <p class="text-xs text-gray-500">Мой профиль</p>
              </div>
            </div>
          </router-link>
        </li>
        <li>
          <button
            @click="handleLogout"
            class="flex items-center gap-3 text-gray-500 hover:text-black transition-colors"
          >
            <svg
              class="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
              />
            </svg>
            <span>Выйти</span>
          </button>
        </li>
      </template>
      <template v-else>
        <li>
          <router-link
            to="/login"
            class="flex items-center gap-3 text-gray-500 hover:text-black transition-colors"
          >
            <svg
              class="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
              />
            </svg>
            <span>Войти</span>
          </router-link>
        </li>
      </template>

      <li>
        <div
          @click="openDrawer"
          class="flex items-center cursor-pointer gap-3 text-gray-500 hover:text-black transition-colors"
        >
          <img src="/cart.svg" alt="Cart" class="filter brightness-90" />
          <b class="text-gray-500">{{ totalPrice }} руб.</b>
        </div>
      </li>

      <li>
        <router-link
          to="/favorites"
          class="flex items-center gap-3 text-gray-500 hover:text-black transition-colors relative"
        >
          <div class="relative">
            <img
              src="/heart.svg"
              alt="Favorites"
              class="filter brightness-90"
            />
            <span
              v-if="favoritesCount > 0"
              class="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center"
            >
              {{ favoritesCount }}
            </span>
          </div>
          <span>Закладки</span>
        </router-link>
      </li>
    </ul>
  </header>
</template>
