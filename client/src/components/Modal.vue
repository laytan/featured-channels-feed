<template>
  <div class="fixed inset-0 z-10 overflow-y-auto">
    <div
      class="flex items-end justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0"
    >
      <div class="fixed inset-0 transition-opacity" aria-hidden="true">
        <div class="absolute inset-0 bg-gray-600 opacity-25"></div>
      </div>

      <span
        class="hidden sm:inline-block sm:align-middle sm:h-screen"
        aria-hidden="true"
        >&#8203;</span
      >
      <div
        ref="modal"
        class="inline-block overflow-hidden text-left align-bottom transition-all transform bg-white rounded-lg shadow-xl sm:my-8 sm:align-middle sm:max-w-lg sm:w-full"
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal"
      >
        <div class="px-4 pt-5 pb-4 bg-white sm:p-6 sm:pb-4">
          <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
            <h3 class="text-2xl leading-6 text-gray-900">Retrieving data</h3>
            <div class="mt-2">
              <p class="text-gray-500">
                Crawling the Youtube website takes some time on a minimal google
                cloud instance. This takes a maximum of 1 minute. Thanks for
                understanding.
              </p>
            </div>
          </div>
        </div>
        <div class="px-4 pb-4 bg-gray-50 sm:px-6">
          <button
            @click="onClose"
            type="button"
            class="px-4 py-2 mt-3 text-base font-medium text-gray-700 bg-white border border-blue-200 rounded-md shadow-sm hover:bg-gray-200 focus:outline-none sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
          >
            Close
          </button>
        </div>
        <div id="loading" class="w-10 h-3 bg-blue-400"></div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent, PropType } from "vue";

export default defineComponent({
  props: {
    onClose: {
      type: Function as PropType<() => void>,
      required: true,
    },
  },
  setup() {
    const modal = ref<null | HTMLElement>(null);

    // Update width css variable
    const setWidthVar = () => {
      document.documentElement.style.setProperty(
        "--modal-width",
        `${modal.value?.clientWidth ?? 1000}px`
      );
    };

    // After mount modal will not be null anymore
    onMounted(() => {
      setWidthVar();
      window.addEventListener("resize", setWidthVar);
    });

    return { modal };
  },
});
</script>

<style scoped>
@keyframes loading {
  0% {
    transform: translateX(-100%);
  }

  100% {
    transform: translateX(calc(100% + var(--modal-width)));
  }
}

#loading {
  animation: loading 1.5s infinite linear;
}
</style>
