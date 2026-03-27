<script setup>
import { ref, onMounted, computed, inject } from "vue";
import { useRouter } from "vue-router";
import CardList from "../components/CardList.vue";

const { favorites, toggleFavorite } = inject("favorites");
const { addToCart } = inject("cart");
const router = useRouter();
const isLoading = ref(false);
const selectedCategory = ref("all");

const categories = computed(() => {
  const cats = new Set(favorites.value.map((item) => item.category));
  return Array.from(cats).sort();
});

const filteredFavorites = computed(() => {
  if (selectedCategory.value === "all") return favorites.value;
  return favorites.value.filter(
    (item) => item.category === selectedCategory.value,
  );
});

onMounted(() => {
  if (!favorites.value) {
    router.push("/login");
  }
});

const removeFavorite = async (product) => {
  await toggleFavorite(product);
};

const removeAllFavorites = async () => {
  if (!confirm("Вы уверены, что хотите удалить все закладки?")) return;
  isLoading.value = true;
  try {
    for (const item of favorites.value) {
      await toggleFavorite(item);
    }
  } catch (error) {
    console.error("Ошибка удаления всех избранных", error);
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="max-w-8xl mx-auto">
    <div class="mb-10">
      <div class="flex justify-between items-start">
        <div>
          <h1 class="text-3xl font-bold text-gray-800">Мои закладки</h1>
          <p class="text-gray-600 mt-2">
            Товары, которые вы сохранили для покупки позже
          </p>
        </div>
        <button
          v-if="favorites?.length"
          @click="removeAllFavorites"
          class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition text-sm"
        >
          Очистить все
        </button>
      </div>
    </div>

    <div v-if="favorites?.length && categories.length" class="mb-6">
      <select
        v-model="selectedCategory"
        class="py-2 px-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="all">Все категории</option>
        <option v-for="cat in categories" :key="cat" :value="cat">
          {{ cat }}
        </option>
      </select>
    </div>

    <div v-if="isLoading" class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"
      ></div>
      <p class="mt-4 text-gray-600">Загрузка...</p>
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
      <p class="text-gray-600 mb-4">
        Показано {{ filteredFavorites.length }} из
        {{ favorites.length }} закладок
      </p>
      <CardList
        :items="
          filteredFavorites.map((item) => ({ ...item, isFavorite: true }))
        "
        is-favorites
        :on-remove="removeFavorite"
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
