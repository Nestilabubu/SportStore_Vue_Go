<script setup>
import { ref, reactive, onMounted, watch } from "vue";
import CardList from "../components/CardList.vue";
import { getProducts } from "../utils/products";
import { inject } from "vue";
import { debounce } from "lodash";

const { addToCart, cart } = inject("cart");
const { favorites, toggleFavorite } = inject("favorites");

const items = ref([]);
const isLoading = ref(false);

const filters = reactive({
  sortBy: "title",
  searchQuery: "",
  category: "all",
  minPrice: "",
  maxPrice: "",
  size: "",
});

const categories = [
  { value: "all", label: "Все категории" },
  { value: "мужской", label: "Мужские" },
  { value: "женский", label: "Женские" },
  { value: "детский", label: "Детские" },
];

const sizes = [
  "XS",
  "S",
  "M",
  "L",
  "XL",
  "XXL",
  "XXXL",
  "110",
  "120",
  "130",
  "140",
  "150",
  "160",
];

const onChangeSelect = (e, filterType) => {
  filters[filterType] = e.target.value;
};

const onChangeSearchInput = debounce((e) => {
  filters.searchQuery = e.target.value;
}, 300);

const resetFilters = () => {
  filters.category = "all";
  filters.minPrice = "";
  filters.maxPrice = "";
  filters.size = "";
  filters.searchQuery = "";
  filters.sortBy = "title";
  fetchItems();
};

const fetchItems = async () => {
  isLoading.value = true;
  try {
    const params = {
      sortBy: filters.sortBy,
      category: filters.category !== "all" ? filters.category : undefined,
      search: filters.searchQuery || undefined,
      minPrice: filters.minPrice || undefined,
      maxPrice: filters.maxPrice || undefined,
      size: filters.size || undefined,
    };
    const data = await getProducts(params);
    items.value = data.map((item) => ({
      ...item,
      isFavorite: favorites.value.some((f) => f.id === item.id),
      isAdded: cart.value.some((c) => c.productId === item.id),
      availableSizes: item.sizes.split(","),
    }));
  } catch (error) {
    console.error("Ошибка загрузки товаров:", error);
  } finally {
    isLoading.value = false;
  }
};

const handleAddToCart = async (product) => {
  await addToCart(product.id, product.selectedSize, product.quantity || 1);
  const idx = items.value.findIndex((i) => i.id === product.id);
  if (idx !== -1) {
    items.value[idx].isAdded = true;
  }
};

const handleToggleFavorite = async (product) => {
  await toggleFavorite(product);
  const idx = items.value.findIndex((i) => i.id === product.id);
  if (idx !== -1) {
    items.value[idx].isFavorite = !items.value[idx].isFavorite;
  }
};

watch(
  [cart, favorites],
  () => {
    items.value = items.value.map((item) => ({
      ...item,
      isFavorite: favorites.value.some((f) => f.id === item.id),
      isAdded: cart.value.some((c) => c.productId === item.id),
    }));
  },
  { deep: true },
);

onMounted(fetchItems);
watch(filters, fetchItems);
</script>

<template>
  <div class="mb-10">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-3xl font-bold">Все спортивные костюмы</h2>
      <div class="flex items-center gap-4">
        <select
          @change="(e) => onChangeSelect(e, 'sortBy')"
          :value="filters.sortBy"
          class="py-2 px-3 border border-gray-200 focus:border-gray-400 rounded-md focus:outline-none"
        >
          <option value="title">По названию</option>
          <option value="price">По цене (дешевые)</option>
          <option value="-price">По цене (дорогие)</option>
        </select>
        <div class="relative">
          <input
            @input="onChangeSearchInput"
            :value="filters.searchQuery"
            type="text"
            class="border border-gray-200 rounded-md py-2 pl-10 pr-4 focus:outline-none focus:border-gray-400"
            placeholder="Поиск по названию..."
          />
          <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
          >
            <img src="/search.svg" alt="Поиск" />
          </div>
        </div>
      </div>
    </div>

    <div class="bg-gray-50 p-6 rounded-xl mb-8">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold">Фильтры</h3>
        <button
          @click="resetFilters"
          class="text-sm text-gray-600 hover:text-gray-800 underline"
        >
          Сбросить фильтры
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2"
            >Категория</label
          >
          <select
            @change="(e) => onChangeSelect(e, 'category')"
            :value="filters.category"
            class="w-full py-2 px-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option
              v-for="cat in categories"
              :key="cat.value"
              :value="cat.value"
            >
              {{ cat.label }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2"
            >Размер</label
          >
          <select
            @change="(e) => onChangeSelect(e, 'size')"
            :value="filters.size"
            class="w-full py-2 px-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Все размеры</option>
            <option v-for="size in sizes" :key="size" :value="size">
              {{ size }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2"
            >Цена от</label
          >
          <input
            v-model.number="filters.minPrice"
            @change="fetchItems"
            type="number"
            min="0"
            placeholder="0"
            class="w-full py-2 px-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2"
            >Цена до</label
          >
          <input
            v-model.number="filters.maxPrice"
            @change="fetchItems"
            type="number"
            min="0"
            placeholder="20000"
            class="w-full py-2 px-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>
    </div>
  </div>

  <div v-if="isLoading" class="text-center py-12">
    <div
      class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"
    ></div>
    <p class="mt-4 text-gray-600">Загрузка товаров...</p>
  </div>

  <div v-else-if="items.length === 0" class="text-center py-12">
    <div class="text-4xl mb-4">🛍️</div>
    <h3 class="text-xl font-semibold mb-2">Товары не найдены</h3>
    <p class="text-gray-600">Попробуйте изменить параметры фильтрации</p>
    <button
      @click="resetFilters"
      class="mt-4 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition"
    >
      Сбросить фильтры
    </button>
  </div>

  <div v-else class="mt-10">
    <p class="text-gray-600 mb-4">Найдено {{ items.length }} товаров</p>
    <CardList
      :items="items"
      @add-to-favorite="handleToggleFavorite"
      @add-to-cart="handleAddToCart"
    />
  </div>
</template>
