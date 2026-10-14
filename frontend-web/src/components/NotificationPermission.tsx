"use client";

import { useEffect } from "react";
import { requestForToken, onMessageListener } from "@/lib/firebase";

export default function NotificationPermission() {
  useEffect(() => {
    // Only run on client
    if (typeof window === "undefined" || !("Notification" in window)) return;

    const setupNotifications = async () => {
      try {
        if (Notification.permission === "granted") {
          const token = await requestForToken();
          if (token) {
            console.log("FCM Token registered.");
          }
        } else if (Notification.permission !== "denied") {
          const permission = await Notification.requestPermission();
          if (permission === "granted") {
            const token = await requestForToken();
            if (token) {
              console.log("FCM Token registered.");
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
      listener.then((payload: unknown) => {
        console.log("Received foreground message: ", payload);
        const msg = payload as { notification?: { title?: string; body?: string } } | null;
        if (msg?.notification?.title) {
          new Notification(msg.notification.title, {
            body: msg.notification.body || "",
            icon: '/firebase-logo.png'
          });
        }
      }).catch(err => console.log('failed: ', err));
    }
  }, []);

  return null; // Silent component, just handles logic
}
