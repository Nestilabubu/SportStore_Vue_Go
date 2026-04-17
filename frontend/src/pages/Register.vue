<script setup>
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { register } from "../utils/auth";
import { inject } from "vue";

const router = useRouter();
const { fetchUser } = inject("auth");

const form = ref({
  fullName: "",
  email: "",
  password: "",
  confirmPassword: "",
  phone: "",
  address: "",
});

const errors = ref({});
const isLoading = ref(false);
const successMessage = ref("");

// Состояние для проверки пароля в реальном времени
const passwordStrength = computed(() => {
  const pwd = form.value.password;
  const checks = {
    length: pwd.length >= 8,
    uppercase: /[A-Z]/.test(pwd),
    digit: /\d/.test(pwd),
    special: /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/.test(pwd),
  };
  const passed = Object.values(checks).filter(Boolean).length;
  return { checks, passed, total: 4 };
});

const isPasswordValid = computed(() => passwordStrength.value.passed === 4);
const isConfirmPasswordValid = computed(
  () =>
    form.value.confirmPassword &&
    form.value.password === form.value.confirmPassword,
);

const passwordStrengthText = computed(() => {
  const pct = (passwordStrength.value.passed / 4) * 100;
  if (pct === 100) return "Отличный пароль";
  if (pct >= 50) return "Средний пароль";
  return "Слабый пароль";
});

const passwordStrengthColor = computed(() => {
  const pct = (passwordStrength.value.passed / 4) * 100;
  if (pct === 100) return "bg-green-500";
  if (pct >= 50) return "bg-yellow-500";
  return "bg-red-500";
});

const validateForm = () => {
  errors.value = {};
  if (!form.value.fullName.trim()) errors.value.fullName = "Введите полное имя";
  if (!form.value.email.trim()) errors.value.email = "Введите email";
  else if (!/^\S+@\S+\.\S+$/.test(form.value.email))
    errors.value.email = "Введите корректный email";

  if (!form.value.password) errors.value.password = "Введите пароль";
  else if (!isPasswordValid.value) {
    errors.value.password =
      "Пароль должен содержать минимум 8 символов, заглавную букву, цифру и спецсимвол";
  }

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

          <!-- Поле пароля с индикатором силы -->
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
              placeholder="Не менее 8 символов"
            />

            <!-- Индикатор силы пароля -->
            <div v-if="form.password" class="mt-2">
              <div class="w-full bg-gray-200 rounded-full h-2">
                <div
                  class="h-2 rounded-full transition-all duration-300"
                  :class="passwordStrengthColor"
                  :style="{ width: (passwordStrength.passed / 4) * 100 + '%' }"
                ></div>
              </div>
              <p class="text-xs mt-1 text-gray-600">
                {{ passwordStrengthText }}
              </p>
              <ul class="text-xs mt-2 space-y-1">
                <li
                  :class="
                    passwordStrength.checks.length
                      ? 'text-green-600'
                      : 'text-gray-500'
                  "
                >
                  ✓ Минимум 8 символов
                </li>
                <li
                  :class="
                    passwordStrength.checks.uppercase
                      ? 'text-green-600'
                      : 'text-gray-500'
                  "
                >
                  ✓ Хотя бы одна заглавная буква
                </li>
                <li
                  :class="
                    passwordStrength.checks.digit
                      ? 'text-green-600'
                      : 'text-gray-500'
                  "
                >
                  ✓ Хотя бы одна цифра
                </li>
                <li
                  :class="
                    passwordStrength.checks.special
                      ? 'text-green-600'
                      : 'text-gray-500'
                  "
                >
                  ✓ Хотя бы один спецсимвол (!@#$%^&* и т.д.)
                </li>
              </ul>
            </div>
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">
              {{ errors.password }}
            </p>
          </div>

          <!-- Подтверждение пароля -->
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
                errors.confirmPassword ||
                (form.confirmPassword && !isConfirmPasswordValid)
                  ? 'border-red-500'
                  : 'border-gray-300'
              "
              placeholder="Повторите пароль"
            />
            <p
              v-if="form.confirmPassword && !isConfirmPasswordValid"
              class="mt-1 text-sm text-red-600"
            >
              Пароли не совпадают
            </p>
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
            :disabled="isLoading || !isPasswordValid || !isConfirmPasswordValid"
            class="w-full py-3 px-4 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
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
