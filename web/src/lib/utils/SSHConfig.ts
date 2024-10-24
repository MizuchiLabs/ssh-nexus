import { pb } from "$lib/client";
import { groups } from "$lib/subscriptions";
import { get } from "svelte/store";

export async function generateConfig() {
  if (!pb.authStore.model) return;
  const request = await pb.send("/api/self/machines", { method: "GET" });
  if (!request.machines.length) return;

  const userKeyName = pb.authStore.model.settings?.ssh_key_name;

  // Base SSH config
  let sshConfig = `Host *\n`;
  sshConfig += `    Protocol 2\n`;
  sshConfig += `    IdentitiesOnly yes\n`;
  sshConfig += `    Compression yes\n`;
  sshConfig += `    ForwardAgent no\n`;
  sshConfig += `    ForwardX11 no\n`;
  sshConfig += `    ForwardX11Trusted no\n`;
  if (userKeyName) {
    sshConfig += `    IdentityFile ~/.ssh/${userKeyName}\n`;
  } else {
    sshConfig += `    IdentityFile ~/.ssh/id_rsa\n\n`;
  }

  for (const machine of request.machines) {
    for (const id of machine.groups) {
      let group = get(groups).find((group) => group.id === id);
      if (!group) continue;
      if (group.linux_username === "root") {
        sshConfig += `Host ${machine.name}\n`;
      } else {
        sshConfig += `Host ${machine.name}-${group.name}\n`;
      }
      sshConfig += `    Hostname ${machine.host}\n`;
      sshConfig += `    User ${group.linux_username}\n`;
      if (machine.port != 22) {
        sshConfig += `    Port ${machine.port}\n`;
      }
      sshConfig += "\n";
    }
  }
  return sshConfig;
}

export function downloadConfig(sshConfig: string) {
  const blob = new Blob([sshConfig], { type: "text/text" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = "config";
  link.click();
  URL.revokeObjectURL(url);
}
