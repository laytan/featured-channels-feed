<template>
  <div class="container px-2 mx-auto mt-5 md:mt-16">
    <div
      v-if="error.length"
      class="fixed top-0 left-0 right-0 m-4 md:left-auto md:m-8"
    >
      <div class="p-6 pr-10 text-red-900 bg-red-200 rounded-lg shadow-md">
        {{ error }}
      </div>
      <button
        @click="error = ''"
        class="absolute top-0 right-0 px-3 py-2 opacity-75 cursor-pointer hover:opacity-100"
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
      <div class="flex flex-wrap p-4 bg-gray-200 shadow-lg md:flex-nowrap">
        <label
          for="channel-search"
          class="mr-4 text-3xl font-semibold text-gray-700"
          >Channel</label
        >
        <div class="flex w-full">
          <input
            id="channel-search"
            v-model="channel"
            class="w-full p-2 border-2 border-blue-200 rounded flex-grow-1"
            type="text"
            placeholder="Try 'LinusTechTips'"
          />
          <button
            class="p-2 pl-4 pr-4 text-white bg-blue-400 rounded hover:bg-blue-300"
          >
            <p class="font-semibold">Go</p>
          </button>
        </div>
      </div>
    </form>
    <div class="mt-2" v-if="result">
      <p class="ml-4 text-sm">
        {{ result.length }} Video{{ result.length > 1 ? "s" : "" }} published by
        featured channels of {{ currChan }} in the last 2 weeks found.
      </p>
      <a
        v-for="video in result"
        :key="video.title"
        :href="`https://www.youtube.com/watch?v=${video.id}`"
        target="_BLANK"
        rel="noopener noreferrer nofollow"
        class="grid grid-rows-1 mx-auto mt-5 md:grid-cols-3 hover:underline"
        style="max-width: 1000px"
      >
        <div class="w-full">
          <img
            class="object-contain w-full"
            :src="video.thumbnail"
            alt="Thumbnail"
          />
        </div>
        <div class="p-3 text-gray-700 bg-gray-200 md:col-span-2">
          <h2
            class="text-2xl font-semibold break-all"
            v-html="sanitize(video.title)"
          ></h2>
          <h3 class="mb-2">
            <span v-html="sanitize(video.channelTitle)"></span> -
            {{ formatDistanceToNow(video.publishedAt) }} ago
          </h3>
          <p class="break-all" v-html="sanitize(video.description)"></p>
        </div>
      </a>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { formatDistanceToNow, compareDesc } from "date-fns";
import DOMPurify from "dompurify";
import Video from "./models/Video";

export default defineComponent({
  name: "App",
  setup() {
    const error = ref("");

    // Modelled to input
    const channel = ref("");

    // Channel we are showing results of
    const currChan = ref("");
    const result = ref<Array<Video> | null>(null);

    const onSearch = async () => {
      try {
        error.value = "";
        const res = await fetch(
          `${import.meta.env.VITE_FUNCTION_URL ?? ""}?channel=${channel.value}`
        );
        const data = await res.json();

        // Only handle responses in the 200 range
        if (!res.ok) {
          error.value =
            data.error ?? "Error retrieving videos, please try again";
          return;
        }

        result.value = data.result
          ?.map(
            (res: any) =>
              new Video(
                res.id ?? "",
                res.channelId ?? "",
                res.channelTitle ?? "",
                res.description ?? "",
                res.publishedAt ?? "",
                res.thumbnail ?? "",
                res.title ?? ""
              )
          )
          .sort((a: Video, b: Video) =>
            compareDesc(a.publishedAt, b.publishedAt)
          );
        currChan.value = channel.value;
      } catch (e) {
        console.error(e);
        error.value = "Error retrieving videos, please try again";
      }
    };

    return {
      onSearch,
      channel,
      error,
      result,
      formatDistanceToNow,
      sanitize: DOMPurify.sanitize,
      currChan,
    };
  },
});
</script>
