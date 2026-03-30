<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { register } from "../utils/auth";
import { inject } from "vue";
const router = useRouter();

const form = ref({
  fullName: "",
  email: "",
  password: "",
  confirmPassword: "",
  phone: "",
  address: "",
});

const { fetchUser } = inject("auth");
const errors = ref({});
const isLoading = ref(false);
const successMessage = ref("");

const validateForm = () => {
  errors.value = {};
  if (!form.value.fullName.trim()) errors.value.fullName = "Введите полное имя";
  if (!form.value.email.trim()) errors.value.email = "Введите email";
  else if (!/^\S+@\S+\.\S+$/.test(form.value.email))
    errors.value.email = "Введите корректный email";
  if (!form.value.password) errors.value.password = "Введите пароль";
  else if (form.value.password.length < 6)
    errors.value.password = "Пароль должен быть не менее 6 символов";
  if (form.value.password !== form.value.confirmPassword)
    errors.value.confirmPassword = "Пароли не совпадают";
  return Object.keys(errors.value).length === 0;
};

const handleSubmit = async () => {
  if (!validateForm()) return;
  isLoading.value = true;
  errors.value = {};
  try {
    await register(form.value);
    await fetchUser();
    successMessage.value = "Регистрация прошла успешно! Перенаправляем...";
    setTimeout(() => router.push("/profile"), 2000);
  } catch (error) {
    if (error.response?.status === 409) {
      errors.value.email = "Пользователь с таким email уже существует";
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
      <h1 class="text-3xl font-bold text-gray-800 mb-2">Регистрация</h1>
      <p class="text-gray-600">
        Создайте учетную запись для доступа к магазину
      </p>
    </div>

    <div
      v-if="successMessage"
      class="mb-6 p-4 bg-green-100 border border-green-400 text-green-700 rounded-lg"
    >
      {{ successMessage }}
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-100">
        <div class="space-y-4">
          <div>
            <label
              for="fullName"
              class="block text-sm font-medium text-gray-700 mb-1"
              >Полное имя *</label
            >
            <input
              id="fullName"
              v-model="form.fullName"
              type="text"
              class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              :class="errors.fullName ? 'border-red-500' : 'border-gray-300'"
              placeholder="Иван Иванов"
            />
            <p v-if="errors.fullName" class="mt-1 text-sm text-red-600">
              {{ errors.fullName }}
            </p>
          </div>

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
              placeholder="Не менее 6 символов"
            />
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">
              {{ errors.password }}
            </p>
          </div>

          <div>
            <label
              for="confirmPassword"
              class="block text-sm font-medium text-gray-700 mb-1"
              >Подтверждение пароля *</label
            >
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              :class="
                errors.confirmPassword ? 'border-red-500' : 'border-gray-300'
              "
              placeholder="Повторите пароль"
            />
            <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">
              {{ errors.confirmPassword }}
            </p>
          </div>

          <div>
            <label
              for="phone"
              class="block text-sm font-medium text-gray-700 mb-1"
              >Телефон</label
            >
            <input
              id="phone"
              v-model="form.phone"
              type="tel"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              placeholder="+7 (999) 999-99-99"
            />
          </div>

          <div>
            <label
              for="address"
              class="block text-sm font-medium text-gray-700 mb-1"
              >Адрес доставки</label
            >
            <textarea
              id="address"
              v-model="form.address"
              rows="3"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              placeholder="Город, улица, дом, квартира"
            ></textarea>
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
            <span v-if="isLoading">Регистрация...</span>
            <span v-else>Зарегистрироваться</span>
          </button>
        </div>

        <div class="mt-6 pt-6 border-t border-gray-200 text-center">
          <p class="text-gray-600">
            Уже есть аккаунт?
            <router-link
              to="/login"
              class="text-blue-600 hover:text-blue-800 font-medium"
              >Войти</router-link
            >
          </p>
        </div>
      </div>
    </form>
  </div>
</template>
