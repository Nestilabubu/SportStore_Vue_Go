<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { inject } from "vue";
import CardList from "../components/CardList.vue";

const { favorites, toggleFavorite } = inject("favorites");
const { addToCart } = inject("cart");
const router = useRouter();
const isLoading = ref(false);

// При монтировании проверяем авторизацию (inject уже содержит favorites, но он мог быть пуст, если не авторизован)
onMounted(() => {
  if (!favorites.value) {
    router.push("/login");
  }
});

const removeFavorite = async (product) => {
  await toggleFavorite(product);
};
</script>

<template>
  <div class="max-w-8xl mx-auto">
    <div class="mb-10">
      <h1 class="text-3xl font-bold text-gray-800">Мои закладки</h1>
      <p class="text-gray-600 mt-2">
        Товары, которые вы сохранили для покупки позже
      </p>
    </div>

    <div v-if="isLoading" class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"
      ></div>
      <p class="mt-4 text-gray-600">Загрузка закладок...</p>
    </div>

    <div v-else-if="!favorites?.length" class="text-center py-12">
      <div class="text-4xl mb-4">⭐</div>
      <h3 class="text-xl font-semibold mb-2">Закладок пока нет</h3>
      <p class="text-gray-600 mb-6">
        Добавляйте понравившиеся товары в закладки
      </p>
      <router-link
        to="/"
        class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
      >
        Перейти к покупкам
      </router-link>
    </div>

    <div v-else>
      <CardList
        :items="favorites.map((item) => ({ ...item, isFavorite: true }))"
        is-favorites
        @add-to-favorite="removeFavorite"
        @add-to-cart="
          (product) =>
            addToCart(
              product.id,
              product.selectedSize || product.sizes.split(',')[0],
            )
        "
      />
    </div>
  </div>
</template>
