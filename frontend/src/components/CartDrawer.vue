<script setup>
import { inject, ref } from "vue";
import { useRouter } from "vue-router";
import { createOrder } from "../utils/orders";

const {
  cart,
  totalPrice,
  vatPrice,
  closeDrawer,
  removeFromCart,
  updateCartQuantity,
} = inject("cart");
const { user } = inject("auth");
const router = useRouter();
const isCreatingOrder = ref(false);

const increaseQuantity = (item) => {
  updateCartQuantity(item.id, (item.quantity || 1) + 1);
};

const decreaseQuantity = (item) => {
  const current = item.quantity || 1;
  if (current > 1) {
    updateCartQuantity(item.id, current - 1);
  } else {
    removeFromCart(item.id);
  }
};

const handleCreateOrder = async () => {
  if (!user.value) {
    alert("Пожалуйста, войдите в систему");
    closeDrawer();
    router.push("/login");
    return;
  }
  isCreatingOrder.value = true;
  try {
    await createOrder();
    alert("Заказ успешно оформлен!");
    closeDrawer();
  } catch (error) {
    console.error("Ошибка при оформлении заказа:", error);
    alert("Произошла ошибка. Попробуйте снова.");
  } finally {
    isCreatingOrder.value = false;
  }
};

const totalItems = () =>
  cart.value.reduce((acc, item) => acc + item.quantity, 0);
</script>

<template>
  <div
    class="fixed top-0 left-0 h-full w-full bg-black z-10 opacity-70"
    @click="closeDrawer"
  ></div>

  <div class="fixed top-0 right-0 h-full w-120 bg-white z-20 p-8 flex flex-col">
    <div class="flex items-center justify-between mb-10">
      <h2 class="text-2xl font-bold">Корзина</h2>
      <button
        @click="closeDrawer"
        class="flex items-center gap-2 text-gray-400 hover:text-black transition"
      >
        <svg
          class="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
        <span class="text-sm">Закрыть</span>
      </button>
    </div>

    <div
      v-if="cart.length === 0"
      class="flex-1 flex flex-col items-center justify-center"
    >
      <div class="text-5xl mb-4">🛒</div>
      <h3 class="text-xl font-semibold mb-2">Корзина пуста</h3>
      <p class="text-gray-500 text-center mb-6">
        Добавьте хотя бы один товар, чтобы сделать заказ
      </p>
      <button
        @click="closeDrawer"
        class="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition"
      >
        Продолжить покупки
      </button>
    </div>

    <div v-else class="flex flex-col h-full">
      <div class="flex-1 overflow-y-auto mb-4">
        <div class="space-y-4">
          <div
            v-for="item in cart"
            :key="item.id"
            class="flex items-center gap-4 border border-gray-200 rounded-xl p-4"
          >
            <img
              :src="item.imageUrl"
              :alt="item.title"
              class="w-20 h-20 object-cover rounded-lg"
            />

            <div class="flex-1">
              <h4 class="font-medium mb-1">{{ item.title }}</h4>
              <div class="flex items-center gap-3 text-sm text-gray-600 mb-2">
                <span>Размер: {{ item.size }}</span>
              </div>

              <div class="flex items-center justify-between">
                <div
                  class="flex items-center border border-gray-300 rounded-lg"
                >
                  <button
                    @click="decreaseQuantity(item)"
                    class="w-8 h-8 flex items-center justify-center text-gray-600 hover:bg-gray-100 rounded-l-lg transition"
                    :class="{
                      'text-gray-300 cursor-not-allowed':
                        (item.quantity || 1) <= 1,
                    }"
                  >
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M20 12H4"
                      />
                    </svg>
                  </button>

                  <span class="w-10 text-center font-medium">{{
                    item.quantity || 1
                  }}</span>

                  <button
                    @click="increaseQuantity(item)"
                    class="w-8 h-8 flex items-center justify-center text-gray-600 hover:bg-gray-100 rounded-r-lg transition"
                  >
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 4v16m8-8H4"
                      />
                    </svg>
                  </button>
                </div>

                <div class="flex items-center gap-4">
                  <span class="font-bold text-right min-w-[100px]">
                    {{
                      (item.price * item.quantity).toLocaleString("ru-RU")
                    }}
                    руб.
                  </span>
                  <button
                    @click="removeFromCart(item.id)"
                    class="text-red-500 hover:text-red-700 transition"
                    title="Удалить"
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
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </button>
                </div>
              </div>

              <div class="text-xs text-gray-500 mt-1">
                {{ item.price.toLocaleString("ru-RU") }} руб. / шт
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="border-t border-gray-200 pt-4">
        <div class="space-y-3">
          <div class="flex justify-between text-gray-600">
            <span>Товары ({{ totalItems() }})</span>
            <span>{{ totalPrice.toLocaleString("ru-RU") }} руб.</span>
          </div>
          <div class="flex justify-between text-gray-600">
            <span>НДС 5%</span>
            <span>{{ vatPrice.toLocaleString("ru-RU") }} руб.</span>
          </div>
          <div
            class="flex justify-between text-xl font-bold pt-3 border-t border-gray-200"
          >
            <span>Итого</span>
            <span
              >{{ (totalPrice + vatPrice).toLocaleString("ru-RU") }} руб.</span
            >
          </div>
        </div>

        <div class="mt-6">
          <button
            @click="handleCreateOrder"
            :disabled="isCreatingOrder"
            class="w-full py-4 bg-green-500 text-white font-bold rounded-xl hover:bg-green-600 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg
              v-if="isCreatingOrder"
              class="animate-spin h-5 w-5 text-white"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            <span v-if="isCreatingOrder">Оформление...</span>
            <span v-else>Оформить заказ →</span>
          </button>
          <p class="text-xs text-gray-500 text-center mt-2">
            Заказ и статистика будут сохранены в вашем профиле
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
}
.overflow-y-auto::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}
.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 3px;
}
.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
