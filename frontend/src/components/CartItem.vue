<script setup>
import { computed } from "vue";

const emit = defineEmits(["onClickRemove"]);

const props = defineProps({
  id: Number,
  title: String,
  imageUrl: String,
  price: Number,
  size: String,
  category: String,
});

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
</script>

<template>
  <div
    class="flex items-center border border-slate-200 p-4 rounded-xl gap-4 hover:bg-gray-50 transition-colors"
  >
    <img
      :src="imageUrl"
      :alt="title"
      class="w-20 h-20 object-cover rounded-lg"
    />

    <div class="flex flex-col flex-1">
      <div class="flex justify-between items-start">
        <div>
          <p class="font-medium text-gray-800">{{ title }}</p>
          <div class="flex items-center gap-2 mt-2 flex-wrap">
            <span
              :class="['text-xs px-2 py-1 rounded font-medium', categoryColor]"
            >
              {{ category }}
            </span>

            <span
              v-if="size"
              class="text-sm px-2 py-1 bg-gray-100 rounded text-gray-700"
            >
              Размер: {{ size }}
            </span>
          </div>
        </div>

        <img
          @click="emit('onClickRemove')"
          class="opacity-60 hover:opacity-100 cursor-pointer transition w-5 h-5"
          src="/close.svg"
          alt="Удалить"
        />
      </div>

      <div class="flex justify-between items-center mt-4">
        <b class="text-lg text-gray-800">{{ price }} руб.</b>
      </div>
    </div>
  </div>
</template>
