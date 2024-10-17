export function StringColor(str: string) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }

  let color = (hash & 0x00ffffff).toString(16).toUpperCase();
  while (color.length < 6) {
    color = "0" + color;
  }
  return `background-color: #${color}`;
}

export function getRandomColor(str: string) {
  const styleColors = [
    "primary",
    "secondary",
    "tertiary",
    "success",
    "warning",
    "error",
    "surface",
  ];
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }

  const colorIndex = Math.abs(hash) % styleColors.length;
  const randomColor = styleColors[colorIndex];
  const randomShade = Math.floor(Math.random() * 5 + 3) * 100;
  return `bg-${randomColor}-${randomShade}`;
}
