<template>
  <div class="max-w-md mx-auto">
    <!-- Сообщение об успешном заказе -->
    <div v-if="orderSuccess" class="text-center">
      <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-100">
        <div class="text-5xl mb-4">🎉</div>
        <h1 class="text-2xl font-bold text-gray-800 mb-3">
          Ваш заказ принят!
        </h1>
        <p class="text-gray-600 mb-6">
          Скоро ваши костюмы к вам приедут.
        </p>
        <p class="text-sm text-gray-500">
          Перенаправление через {{ countdown }} секунд...
        </p>
      </div>
    </div>

    <!-- Форма оплаты -->
    <div v-else>
      <div class="text-center mb-10">
        <h1 class="text-3xl font-bold text-gray-800 mb-2">Оплата заказа</h1>
        <p class="text-gray-600">Введите данные карты для оплаты</p>
      </div>

      <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-100">
        <div class="mb-6 p-4 bg-gray-50 rounded-lg text-center">
          <p class="text-gray-600">Сумма к оплате:</p>
          <p class="text-2xl font-bold text-gray-800">{{ totalAmount }} руб.</p>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <!-- Поле номера карты с иконкой -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Номер карты</label>
            <div class="relative">
              <input
                v-model="cardNumber"
                type="text"
                placeholder="1234 5678 9012 3456"
                maxlength="19"
                @input="formatCardNumber"
                class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                :class="errors.cardNumber ? 'border-red-500' : 'border-gray-300'"
              />
              <div v-if="cardType" class="absolute right-3 top-1/2 transform -translate-y-1/2">
                <span v-if="cardType === 'visa'" class="px-2 py-1 text-xs font-bold text-white bg-blue-600 rounded">VISA</span>
                <span v-else-if="cardType === 'mastercard'" class="px-2 py-1 text-xs font-bold text-white bg-red-600 rounded">MC</span>
                <span v-else-if="cardType === 'mir'" class="px-2 py-1 text-xs font-bold text-white bg-green-600 rounded">МИР</span>
              </div>
            </div>
            <p v-if="errors.cardNumber" class="mt-1 text-sm text-red-600">{{ errors.cardNumber }}</p>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Срок (ММ/ГГ)</label>
              <input
                v-model="cardExpiry"
                type="text"
                placeholder="MM/YY"
                maxlength="5"
                @input="formatExpiry"
                class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                :class="errors.cardExpiry ? 'border-red-500' : 'border-gray-300'"
              />
              <p v-if="errors.cardExpiry" class="mt-1 text-sm text-red-600">{{ errors.cardExpiry }}</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">CVV</label>
              <input
                v-model="cardCvv"
                type="password"
                placeholder="123"
                maxlength="3"
                class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                :class="errors.cardCvv ? 'border-red-500' : 'border-gray-300'"
              />
              <p v-if="errors.cardCvv" class="mt-1 text-sm text-red-600">{{ errors.cardCvv }}</p>
            </div>
          </div>

          <div v-if="submitError" class="p-3 bg-red-50 border border-red-200 rounded-lg">
            <p class="text-sm text-red-600">{{ submitError }}</p>
          </div>

          <button
            type="submit"
            :disabled="isProcessing"
            class="w-full py-3 px-4 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition disabled:opacity-50"
          >
            <span v-if="isProcessing">Обработка...</span>
            <span v-else>Оплатить {{ totalAmount }} руб.</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, inject } from 'vue';
import { useRouter } from 'vue-router';
import { createOrder } from '../utils/orders';
import { getCart } from '../utils/cart';

const router = useRouter();
const { closeDrawer, refreshCart } = inject('cart');
const { user } = inject('auth');

const cardNumber = ref('');
const cardExpiry = ref('');
const cardCvv = ref('');
const isProcessing = ref(false);
const submitError = ref('');
const orderSuccess = ref(false);
const countdown = ref(5);
const errors = ref({
  cardNumber: '',
  cardExpiry: '',
  cardCvv: ''
});

const cart = ref([]);
const totalPrice = ref(0);
const vatPrice = ref(0);

const totalAmount = computed(() => totalPrice.value + vatPrice.value);

const cardType = computed(() => {
  const number = cardNumber.value.replace(/\s/g, '');
  if (!number) return null;
  if (number[0] === '4') return 'visa';
  if (number.length >= 2 && number[0] === '5' && '12345'.includes(number[1])) return 'mastercard';
  if (number.length >= 4 && number.startsWith('220') && '01234'.includes(number[3])) return 'mir';
  return null;
});

const loadCart = async () => {
  try {
    const cartData = await getCart();
    cart.value = cartData;
    totalPrice.value = cartData.reduce((acc, item) => acc + item.price * item.quantity, 0);
    vatPrice.value = Math.round(totalPrice.value * 0.05);
  } catch (error) {
    console.error('Ошибка загрузки корзины', error);
  }
};

const formatCardNumber = () => {
  let value = cardNumber.value.replace(/\s/g, '');
  if (value.length > 16) value = value.slice(0, 16);
  const parts = value.match(/.{1,4}/g);
  if (parts) {
    cardNumber.value = parts.join(' ');
  } else {
    cardNumber.value = value;
  }
};

const formatExpiry = () => {
  let value = cardExpiry.value.replace(/\D/g, '');
  if (value.length > 4) value = value.slice(0, 4);
  if (value.length >= 3) {
    cardExpiry.value = `${value.slice(0, 2)}/${value.slice(2)}`;
  } else {
    cardExpiry.value = value;
  }
};

const validate = () => {
  let valid = true;
  errors.value = { cardNumber: '', cardExpiry: '', cardCvv: '' };

  const cardNumDigits = cardNumber.value.replace(/\s/g, '');
  if (!/^\d{16}$/.test(cardNumDigits)) {
    errors.value.cardNumber = 'Введите 16 цифр номера карты';
    valid = false;
  }

  const expiryMatch = cardExpiry.value.match(/^(\d{2})\/(\d{2})$/);
  if (!expiryMatch) {
    errors.value.cardExpiry = 'Используйте формат ММ/ГГ';
    valid = false;
  } else {
    const month = parseInt(expiryMatch[1], 10);
    const year = parseInt(expiryMatch[2], 10);
    const now = new Date();
    const currentYear = now.getFullYear() % 100;
    const currentMonth = now.getMonth() + 1;
    if (month < 1 || month > 12) {
      errors.value.cardExpiry = 'Неверный месяц';
      valid = false;
    } else if (year < currentYear || (year === currentYear && month < currentMonth)) {
      errors.value.cardExpiry = 'Карта просрочена';
      valid = false;
    }
  }

  if (!/^\d{3}$/.test(cardCvv.value)) {
    errors.value.cardCvv = 'Введите 3 цифры CVV';
    valid = false;
  }

  return valid;
};

const handleSubmit = async () => {
  if (!validate()) return;
  if (!user.value) {
    router.push('/login');
    return;
  }

  isProcessing.value = true;
  submitError.value = '';

  try {
    await createOrder();
    await refreshCart();
    orderSuccess.value = true;

    const timer = setInterval(() => {
      if (countdown.value <= 1) {
        clearInterval(timer);
        router.push('/');
      } else {
        countdown.value--;
      }
    }, 1000);
  } catch (error) {
    console.error('Ошибка оплаты', error);
    submitError.value = 'Ошибка при оформлении заказа. Попробуйте снова.';
    isProcessing.value = false;
  }
};

onMounted(async () => {
  closeDrawer();
  if (!user.value) {
    router.push('/login');
    return;
  }
  await loadCart();
  if (cart.value.length === 0) {
    router.push('/');
  }
});
</script>