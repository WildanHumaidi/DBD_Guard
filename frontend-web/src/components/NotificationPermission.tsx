"use client";

import { useEffect, useState } from "react";
import { requestForToken, onMessageListener } from "@/lib/firebase";
import api from "@/lib/api";

export default function NotificationPermission() {
  const [notificationStatus, setNotificationStatus] = useState<string>("default");

  useEffect(() => {
    // Only run on client
    if (typeof window === "undefined" || !("Notification" in window)) return;

    setNotificationStatus(Notification.permission);

    const setupNotifications = async () => {
      try {
        if (Notification.permission === "granted") {
          const token = await requestForToken();
          if (token) {
            console.log("FCM Token registered.");
            // Optionally, send this token to backend to update user profile
            // await api.post("/users/fcm-token", { token });
          }
        } else if (Notification.permission !== "denied") {
          const permission = await Notification.requestPermission();
          setNotificationStatus(permission);
          if (permission === "granted") {
            const token = await requestForToken();
            if (token) {
              console.log("FCM Token registered.");
              // await api.post("/users/fcm-token", { token });
            }
          }
        }
      } catch (e) {
        console.error("Error setting up notifications", e);
      }
    };

    setupNotifications();

    // Listen for foreground messages
    const listener = onMessageListener();
    if (listener) {
      listener.then((payload: any) => {
        console.log("Received foreground message: ", payload);
        // Could show a toast notification here
        if (payload?.notification) {
          new Notification(payload.notification.title, {
            body: payload.notification.body,
            icon: '/firebase-logo.png' // optional icon
          });
        }
      }).catch(err => console.log('failed: ', err));
    }
  }, []);

  return null; // Silent component, just handles logic
}
