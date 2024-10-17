import type {
  ClientResponseError,
  RecordModel,
  RecordSubscription,
} from "pocketbase";
import { pb } from "./client";
import { writable } from "svelte/store";

export const users = writable([] as RecordModel[]);
export const groups = writable([] as RecordModel[]);
export const machines = writable([] as RecordModel[]);
export const permissions = writable([] as RecordModel[]);
export const providers = writable([] as RecordModel[]);
export const settings = writable([] as RecordModel[]);
export const requests = writable([] as RecordModel[]);
export const tags = writable([] as RecordModel[]);
export const auditlog = writable([] as RecordModel[]);

export const collections = [
  "users",
  "groups",
  "machines",
  "permissions",
  "providers",
  "settings",
  "requests",
  "tags",
  "auditlog",
];

export const subscribeRecords = async () => {
  collections.forEach(async (collection) => {
    try {
      // Fetch initial data
      const store = getStore(collection);
      store?.set(
        await pb.collection(collection).getFullList({
          expand: getExpands(collection),
          requestKey: collection,
        }),
      );

      // Subscribe to SSE
      pb.realtime.subscribe(
        collection,
        async (event: RecordSubscription<RecordModel>) => {
          await updateCollection(collection, event);
        },
        {
          expand: getExpands(collection),
        },
      );
    } catch (error: ClientResponseError | any) {
      console.error(
        `Failed to subscribe to collection ${collection}, error: `,
        error.message,
      );
    }
  });
};

// Return the correct store based on collection name
const getStore = (collection: string) => {
  switch (collection) {
    case "users":
      return users;
    case "groups":
      return groups;
    case "machines":
      return machines;
    case "permissions":
      return permissions;
    case "providers":
      return providers;
    case "settings":
      return settings;
    case "requests":
      return requests;
    case "tags":
      return tags;
    case "auditlog":
      return auditlog;
    default:
      return undefined;
  }
};

// Depending on collection we want relational data expanded
const getExpands = (collection: string) => {
  switch (collection) {
    case "users":
      return "permission,groups";
    case "machines":
      return "users,groups,tags,provider";
    case "permissions":
      return "users,groups,machines";
    case "requests":
      return "user,group,machine";
    case "auditlog":
      return "user";
    default:
      return "";
  }
};

// Updates entries
export const updateCollection = async (
  collection: string,
  event: RecordSubscription<RecordModel>,
) => {
  const action = event.action;
  const record = event.record;
  const store = getStore(collection);
  if (!store) {
    console.error(`Unknown collection: ${collection}`);
    return;
  }

  store.update((data) => {
    const index = data.findIndex((d) => d.id === record.id);
    switch (action) {
      case "create":
        if (index === -1) {
          data.push(record);
        }
        break;
      case "update":
        if (index !== -1) {
          data[index] = record;
        }
        break;
      case "delete":
        if (index !== -1) {
          data.splice(index, 1);
        }
        break;
      default:
        console.error(`Unknown action: ${action}`);
    }
    return data;
  });
};

// For filtering
export const filterCollection = async (
  collection: string,
  page = 1,
  perPage = 30,
  filter = "",
  sort = "-created",
) => {
  const store = getStore(collection);
  const expands = getExpands(collection);

  try {
    const records = await pb.collection(collection).getList(page, perPage, {
      sort: sort,
      filter: filter,
      expand: expands,
      requestKey: null,
    });

    store?.set(records.items);

    return records.totalItems;
  } catch (error: ClientResponseError | any) {
    console.error(`Failed to get ${collection}, error: `, error.message);
    return 0;
  }
};

export const unsubscribeRecords = () => {
  try {
    pb.realtime.unsubscribe();
    users.set([]);
    groups.set([]);
    machines.set([]);
    permissions.set([]);
    providers.set([]);
    settings.set([]);
    requests.set([]);
    tags.set([]);
    auditlog.set([]);
  } catch (error: ClientResponseError | any) {
    console.error(`Failed to unsubscribe from records, error: `, error.message);
  }
};
