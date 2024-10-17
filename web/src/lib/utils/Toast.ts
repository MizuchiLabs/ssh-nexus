import type { ToastSettings, ToastStore } from "@skeletonlabs/skeleton";

const toastBackgrounds: Record<string, string> = {
  error: "variant-filled-error",
  warning: "variant-filled-warning",
  success: "variant-filled-success",
};

export const showToast = (
  store: ToastStore,
  message: string,
  type = "success",
  options?: Partial<ToastSettings>,
) => {
  const background = toastBackgrounds[type];

  const toastSettings: ToastSettings = {
    message,
    background,
    autohide: true,
    timeout: 3000,
    ...options, // Override default options with provided options
  };

  store.trigger(toastSettings);
};
