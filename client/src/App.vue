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
            :disabled="loading"
          />
          <button
            :class="`p-2 pl-4 pr-4 flex items-center text-white rounded ${
              loading ? 'bg-gray-400' : 'bg-blue-400 hover:bg-blue-300'
            }`"
          >
            <p class="font-semibold">Go</p>
            <svg
              v-if="loading"
              class="w-5 h-5 ml-3 text-white animate-spin"
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
          </button>
        </div>
      </div>
    </form>
    <div class="mt-2" v-if="videos.length">
      <p class="ml-4 text-sm">
        Showing {{ videos.length }} video{{
          videos.length > 1 ? "s" : ""
        }}
        published by featured channels of {{ currChan }} (the 8 latest videos
        per channel).
      </p>
      <a
        v-for="video in videos"
        :key="video.title"
        :href="`https://www.youtube.com${video.url}`"
        target="_BLANK"
        rel="noopener noreferrer nofollow"
        class="grid grid-rows-1 mx-auto mt-5 md:grid-cols-3 hover:underline"
        style="max-width: 1000px"
      >
        <div class="w-full">
          <img
            class="object-contain w-full"
            :src="video.thumbnailUrl"
            alt="Thumbnail"
            loading="lazy"
          />
        </div>
        <div class="p-3 text-gray-700 bg-gray-200 md:col-span-2">
          <h2 class="text-2xl font-semibold break-all">
            {{ video.title }}
          </h2>
          <h3 class="mb-2">
            {{ video.channelTitle }} - {{ video.publishedAt }} -
            {{ video.views }}
          </h3>
        </div>
      </a>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, HTMLAttributes } from "vue";
import Channel from "./models/Channel";
import Video from "./models/Video";

export default defineComponent({
  name: "App",
  setup() {
    const error = ref("");
    const loading = ref(false);

    // Modelled to input
    const channel = ref("");

    // Channel we are showing results of
    const currChan = ref("");
    const result = ref<Array<Channel>>([]);
    const videos = computed(() =>
      result.value
        .reduce((acc: Array<Video>, curr) => [...acc, ...curr.latestVideos], [])
        .sort((a, b) => a.getSeconds() - b.getSeconds())
    );

    // Gets and parses results from our cloud function that scrapes youtube
    const onSearch = async () => {
      if (loading.value) {
        return;
      }

      try {
        loading.value = true;
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

        result.value =
          data.channels?.map(
            (res: any): Channel =>
              new Channel(
                res.subscriberCount ?? "",
                res.urlTitle ?? "",
                res.displayTitle ?? "",
                res.latestVideos?.map(
                  (video: any): Video =>
                    new Video(
                      video.url ?? "",
                      video.publishedAt ?? "",
                      video.thumbnail ?? "",
                      video.title ?? "",
                      video.views ?? "",
                      res.displayTitle ?? ""
                    )
                ) ?? []
              )
          ) ?? [];
        currChan.value = channel.value;
        channel.value = "";
      } catch (e) {
        console.error(e);
        error.value = "Error retrieving videos, please try again";
      } finally {
        loading.value = false;
      }
    };

    return {
      onSearch,
      channel,
      error,
      videos,
      currChan,
      loading,
    };
  },
});
</script>
