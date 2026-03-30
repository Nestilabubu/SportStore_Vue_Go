<script setup>
import { ref, computed, watch, onMounted, inject } from "vue";

const props = defineProps({
  id: Number,
  imgUrl: String,
  title: String,
  price: Number,
  isFavorite: Boolean, // больше не используем, оставлен для совместимости
  isAdded: Boolean,
  onClickAdd: Function,
  onClickFav: Function,
  onRemove: Function,
  category: String,
  availableSizes: Array,
});

const emit = defineEmits(["addToFavorite", "addToCart"]);

const cart = inject("cart");
const favorites = inject("favorites"); // получаем реактивный ref из App.vue

// Вычисляем, находится ли товар в избранном
const isFavoriteLocal = computed(() => {
  return favorites?.favorites?.value?.some((f) => f.id === props.id) ?? false;
});

const selectedSize = ref(props.availableSizes?.[0] || "");

// Вычисляем элемент корзины для данного товара и выбранного размера
const cartItem = computed(() => {
  if (!cart?.cart?.value) return null;
  return cart.cart.value.find(
    (item) => item.productId === props.id && item.size === selectedSize.value,
  );
});

const cartQuantity = computed(() => cartItem.value?.quantity || 0);
const isInCart = computed(() => cartQuantity.value > 0);

// Сохраняем выбранный размер в localStorage
watch(selectedSize, (newSize) => {
  if (props.id) {
    const savedSizes = JSON.parse(
      localStorage.getItem("selectedSizes") || "{}",
    );
    savedSizes[props.id] = newSize;
    localStorage.setItem("selectedSizes", JSON.stringify(savedSizes));
  }
});

const savedSizes = JSON.parse(localStorage.getItem("selectedSizes") || "{}");
if (
  props.id &&
  savedSizes[props.id] &&
  props.availableSizes?.includes(savedSizes[props.id])
) {
  selectedSize.value = savedSizes[props.id];
}

const categoryColor = computed(() => {
  switch (props.category) {
    case "мужской":
      return "bg-blue-100 text-blue-800";
    case "женский":
      return "bg-pink-100 text-pink-800";
    case "детский":
      return "bg-green-100 text-green-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
});

const handleFavoriteClick = () => {
  if (props.onClickFav) {
    // Вызываем родительский обработчик, который обновит глобальное состояние
    props.onClickFav();
    // Ничего больше не делаем, computed сам обновится при изменении favorites
  }
};

// Добавление в корзину
const addToCartHandler = () => {
  if (props.onClickAdd) {
    const itemData = {
      id: props.id,
      title: props.title,
      price: props.price,
      imageUrl: props.imgUrl,
      category: props.category,
      selectedSize: selectedSize.value,
      quantity: 1,
    };
    props.onClickAdd(itemData);
  }
};

const incrementQuantity = async () => {
  if (cartItem.value) {
    await cart.updateCartQuantity?.(cartItem.value.id, cartQuantity.value + 1);
  } else {
    addToCartHandler();
  }
};

const decrementQuantity = async () => {
  if (cartItem.value) {
    const newQuantity = cartQuantity.value - 1;
    if (newQuantity <= 0) {
      await cart.removeFromCart?.(cartItem.value.id);
    } else {
      await cart.updateCartQuantity?.(cartItem.value.id, newQuantity);
    }
  }
};

const handleRemove = () => {
  if (props.onRemove) {
    props.onRemove({
      id: props.id,
      title: props.title,
      price: props.price,
      imageUrl: props.imgUrl,
      category: props.category,
    });
  }
};

// Не нужны onMounted и обработчики storage, так как используем глобальное состояние
</script>

<template>
  <div
    class="relative bg-white border border-slate-100 rounded-3xl p-8 cursor-pointer transition hover:-translate-y-2 hover:shadow-xl"
  >
    <div class="absolute top-3 left-3 z-10">
      <span
        :class="['px-3 py-1 text-sm font-medium rounded-full', categoryColor]"
      >
        {{ category }}
      </span>
    </div>

    <img
      v-if="onClickFav"
      @click="handleFavoriteClick"
      :src="isFavoriteLocal ? '/like-2.svg' : '/like-1.svg'"
      alt="Добавить в избранное"
      class="absolute top-8 right-8 z-10 w-6 h-6 cursor-pointer hover:scale-110 transition"
    />

    <div
      class="rounded-2xl overflow-hidden mb-4 h-64 flex items-center justify-center bg-gray-100"
    >
      <img
        :src="imgUrl"
        :alt="title"
        class="w-full h-full object-cover transition-transform duration-300 hover:scale-105"
      />
    </div>

    <h3 class="mt-2 font-medium text-lg text-gray-800 line-clamp-2 h-14">
      {{ title }}
    </h3>

    <div class="mt-3">
      <p class="text-sm text-slate-500 mb-1">Размер:</p>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="size in availableSizes"
          :key="size"
          @click.stop="selectedSize = size"
          class="px-3 py-1 text-sm border rounded-md transition min-w-[40px]"
          :class="
            selectedSize === size
              ? 'border-blue-500 bg-blue-50 text-blue-600'
              : 'border-gray-200 hover:border-gray-300 text-gray-700'
          "
        >
          {{ size }}
        </button>
      </div>
    </div>

    <div class="flex justify-between items-center mt-5 pt-4 border-t">
      <div class="flex flex-col">
        <span class="text-slate-400 text-sm">Цена:</span>
        <b class="text-xl text-gray-800">{{ price }} руб.</b>
      </div>
      <div class="flex gap-2 items-center">
        <!-- Если товар уже в корзине, показываем контролы количества -->
        <template v-if="isInCart">
          <button
            @click.stop="decrementQuantity"
            class="w-8 h-8 flex items-center justify-center border border-gray-300 rounded-lg hover:bg-gray-100 transition"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="w-5 h-5 text-gray-600"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M20 12H4"
              />
            </svg>
          </button>
          <span class="w-8 text-center font-medium">{{ cartQuantity }}</span>
          <button
            @click.stop="incrementQuantity"
            class="w-8 h-8 flex items-center justify-center border border-gray-300 rounded-lg hover:bg-gray-100 transition"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="w-5 h-5 text-gray-600"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 4v16m8-8H4"
              />
            </svg>
          </button>
        </template>
        <!-- Иначе – кнопка добавления -->
        <img
          v-else-if="onClickAdd"
          @click.stop="addToCartHandler"
          src="/plus.svg"
          alt="Добавить в корзину"
          class="w-8 h-8 cursor-pointer hover:scale-110 transition"
        />
        <!-- Кнопка удаления из избранного -->
        <svg
          v-if="onRemove"
          @click="handleRemove"
          xmlns="http://www.w3.org/2000/svg"
          class="w-7 h-7 cursor-pointer hover:scale-110 transition text-red-500"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M18 6L6 18M6 6l12 12" />
        </svg>
      </div>
    </div>
  </div>
</template>
