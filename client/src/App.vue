<template>
  <div class="container mx-auto mt-5 md:mt-16 px-2">
    <div
      v-if="error.length"
      class="fixed top-0 left-0 right-0 m-4 md:left-auto md:m-8"
    >
      <div class="bg-red-200 text-red-900 rounded-lg shadow-md p-6 pr-10">
        {{ error }}
      </div>
      <button
        @click="error = ''"
        class="opacity-75 cursor-pointer absolute top-0 right-0 py-2 px-3 hover:opacity-100"
      >
        x
      </button>
    </div>

    <div class="pl-4 mb-5">
      <h1 class="text-3xl font-semibold text-gray-700">
        Featured Channel Feed
      </h1>
      <p class="text-2xl text-gray-700">
        All Featured Channels in
        <strong class="text-blue-400">one place</strong>.
      </p>
    </div>
    <form @submit.prevent="onSearch">
      <div class="bg-gray-200 shadow-lg p-4 flex flex-wrap md:flex-nowrap">
        <label
          for="channel-search"
          class="text-3xl text-gray-700 mr-4 font-semibold"
          >Channel</label
        >
        <div class="flex w-full">
          <input
            id="channel-search"
            v-model="channel"
            class="w-full rounded p-2 flex-grow-1 border-2 border-blue-200"
            type="text"
            placeholder="Try 'LinusTechTips'"
          />
          <button
            class="bg-blue-400 hover:bg-blue-300 rounded text-white p-2 pl-4 pr-4"
          >
            <p class="font-semibold">Go</p>
          </button>
        </div>
      </div>
    </form>
    <a
      href="https://www.youtube.com/watch?v=CDAsYgAXeb0"
      class="grid grid-rows-1 md:grid-cols-3 mx-auto mt-5 hover:underline"
      style="max-width: 1000px"
    >
      <div class="relative w-full">
        <img
          class="w-full"
          src="https://i.ytimg.com/vi/CDAsYgAXeb0/hqdefault.jpg?sqp=-oaymwEZCPYBEIoBSFXyq4qpAwsIARUAAIhCGAFwAQ==&rs=AOn4CLAYVFUjIravcWsr8t6t-fsDs2Em7Q"
          alt="Thumbnail"
        />
        <span
          class="absolute bottom-0 right-0 mr-2 mb-1 p-1 bg-black text-gray-100"
          >4:03</span
        >
      </div>
      <div class="md:col-span-2 text-gray-700 bg-gray-200 p-3">
        <h2 class="text-2xl font-semibold">
          Lijpe - Mansory ft. Frenna (prod. Trobi & Vanno)
        </h2>
        <h3 class="mb-2">TopNotch - 429k weergaven - 3 dagen geleden</h3>
        <p>
          Stream of download 'Mansory': http://Lijpe.lnk.to/Mansory Offficial
          Music Video Lijpe - Mansory ft. Frenna (prod. Trobi & Vanno) Video
          Credits: Director: Caio Silva Producent: Rabia Abdoella...
        </p>
      </div>
    </a>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";

export default defineComponent({
  name: "App",
  setup() {
    const channel = ref("");
    const error = ref("");
    const onSearch = async () => {
      try {
        error.value = "";

        const res = await fetch(
          `http://localhost:8080?channel=${channel.value}`
        );

        // Only handle responses in the 200 range
        if (!res.ok) {
          throw await res.json();
        }

        const data = await res.json();
        console.log(data);
      } catch (e) {
        if (e.error?.length) {
          error.value = e.error;
        } else if (e.length) {
          error.value = e;
        } else if (e.error?.message?.length) {
          error.value = e.error.message;
        } else {
          error.value = e.toString();
        }
      }
    };

    return { onSearch, channel, error };
  },
});
</script>
