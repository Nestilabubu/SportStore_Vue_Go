<script setup>
import { inject } from "vue";
import CartItem from "./CartItem.vue";

const { cart, removeFromCart } = inject(["cart"]);
</script>

<template>
  <div v-auto-animate class="flex flex-col flex-1 gap-4">
    <CartItem
      v-for="item in cart"
      :key="item.id"
      :id="item.id"
      :title="item.title"
      :price="item.price"
      :image-url="item.imageUrl"
      :size="item.size || item.selectedSize"
      :category="item.category"
      @on-click-remove="() => removeFromCart(item)"
    />
  </div>
</template>
