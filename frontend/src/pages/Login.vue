<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { login } from "../utils/auth";
import { inject } from "vue";

const router = useRouter();
const form = ref({ email: "", password: "" });
const errors = ref({});
const isLoading = ref(false);
const { fetchUser } = inject("auth");

const validateForm = () => {
  errors.value = {};
  if (!form.value.email.trim()) errors.value.email = "Введите email";
  else if (!/^\S+@\S+\.\S+$/.test(form.value.email))
    errors.value.email = "Введите корректный email";
  if (!form.value.password) errors.value.password = "Введите пароль";
  return Object.keys(errors.value).length === 0;
};

const handleSubmit = async () => {
  if (!validateForm()) return;
  isLoading.value = true;
  errors.value = {};
  try {
    await login(form.value);
    await fetchUser(); // обновляем состояние
    router.push("/profile");
  } catch (error) {
    if (error.response?.status === 401) {
      errors.value.submit = "Неверный email или пароль";
    } else {
      errors.value.submit = "Ошибка сервера. Попробуйте позже.";
    }
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="max-w-md mx-auto">
    <div class="text-center mb-10">
      <h1 class="text-3xl font-bold text-gray-800 mb-2">Вход в аккаунт</h1>
      <p class="text-gray-600">Войдите для доступа к вашему профилю</p>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-100">
        <div class="space-y-4">
          <div>
            <label
              for="email"
              class="block text-sm font-medium text-gray-700 mb-1"
              >Email *</label
            >
            <input
              id="email"
              v-model="form.email"
              type="email"
              class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              :class="errors.email ? 'border-red-500' : 'border-gray-300'"
              placeholder="example@mail.com"
            />
            <p v-if="errors.email" class="mt-1 text-sm text-red-600">
              {{ errors.email }}
            </p>
          </div>

          <div>
            <label
              for="password"
              class="block text-sm font-medium text-gray-700 mb-1"
              >Пароль *</label
            >
            <input
              id="password"
              v-model="form.password"
              type="password"
              class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              :class="errors.password ? 'border-red-500' : 'border-gray-300'"
              placeholder="Введите пароль"
            />
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">
              {{ errors.password }}
            </p>
          </div>

          <div
            v-if="errors.submit"
            class="p-3 bg-red-50 border border-red-200 rounded-lg"
          >
            <p class="text-sm text-red-600">{{ errors.submit }}</p>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-3 px-4 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition disabled:opacity-50"
          >
            <span v-if="isLoading">Вход...</span>
            <span v-else>Войти</span>
          </button>
        </div>

        <div class="mt-6 pt-6 border-t border-gray-200 space-y-3 text-center">
          <p class="text-gray-600">
            Нет аккаунта?
            <router-link
              to="/register"
              class="text-blue-600 hover:text-blue-800 font-medium"
            >
              Зарегистрироваться
            </router-link>
          </p>
          <router-link
            to="/"
            class="text-gray-600 hover:text-gray-800 text-sm block"
          >
            ← Вернуться на главную
          </router-link>
        </div>
      </div>
    </form>
  </div>
</template>
