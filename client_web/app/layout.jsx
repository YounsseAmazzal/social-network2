import "./globals.css";
import { AuthSessionProvider } from "@/providers/AuthSessionProvider";
import { UnreadConversationProvider } from "@/providers/UnreadConversationProvider";
import { NotificationCenterProvider } from "@/providers/NotificationCenterProvider";
import { RealtimeSocketProvider } from "@/providers/RealtimeSocketProvider";
import { ToastMessageProvider } from "@/providers/ToastMessageProvider";
import ApplicationShell from "@/app/components/ApplicationShell";
export const metadata = {
  title: {
    default: "Solo Social Network",
    template: "%s | Social Network",
  },
  description:
    "Connect with people, share posts, and join group conversations.",
};
export default function AppRouterLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <ToastMessageProvider>
          <AuthSessionProvider>
            <RealtimeSocketProvider>
              <NotificationCenterProvider>
                <UnreadConversationProvider>
                  <ApplicationShell>{children}</ApplicationShell>
                </UnreadConversationProvider>
              </NotificationCenterProvider>
            </RealtimeSocketProvider>
          </AuthSessionProvider>
        </ToastMessageProvider>
      </body>
    </html>
  );
}
