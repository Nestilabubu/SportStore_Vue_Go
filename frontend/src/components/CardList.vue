<script setup>
import Card from "./Card.vue";

defineProps({
  items: Array,
  isFavorites: Boolean,
});

const emit = defineEmits(["addToFavorite", "addToCart"]);

const handleAddToFavorite = (item) => {
  emit("addToFavorite", item);
};
</script>

<template>
  <div
    v-auto-animate
    class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5"
  >
    <Card
      v-for="item in items"
      :key="item.id"
      :id="item.id"
      :title="item.title"
      :price="item.price"
      :img-url="item.imageUrl"
      :category="item.category"
      :available-sizes="item.availableSizes || []"
      :on-click-add="isFavorites ? null : () => emit('addToCart', item)"
      :on-click-fav="
        isFavorites
          ? null
          : () =>
              handleAddToFavorite({
                ...item,
                selectedSize:
                  item.selectedSize || item.availableSizes?.[0] || '',
              })
      "
      :is-favorite="item.isFavorite"
      :is-added="item.isAdded"
    />
  </div>
</template>
