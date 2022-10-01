import Head from "next/head";
import Image from "next/image";
import React, { DOMElement, useState, useEffect } from "react";
import Divider from "../components/divider";
import SlackButton from "../components/slackButton";
import { useSession, signIn, signOut } from "next-auth/react";
import { getUserSettings } from "../lib/apiClient";
import Link from "next/link";

export default function Dashboard() {
  const [userSettings, setUserSettings] = useState(null);
  const [isLoading, setLoading] = useState(true);

  useEffect(() => {
    getUserSettings().then((data) => {
      setUserSettings(data);
      setLoading(false);
    });
  }, []);

  if (isLoading) return <p>Loading...</p>;
  if (!userSettings) return <p>Failed to load userSettings</p>;
  return (
    <div className="container mx-auto">
      <p className="font-sans text-2xl"> 👋 Hi there</p>
      <Divider />
      <p>Slack Integration</p>
      <p>
        Connection Status:{" "}
        {userSettings.slackConnected ? (
          <p className="text-green-600 font-bold"> Connected</p>
        ) : (
          <p className="text-red-600 font-bold"> Not connected</p>
        )}
      </p>
      <SlackButton
        title="Connect to"
        redirectURI={userSettings.slackRedirectURI}
        userId={userSettings.userId}
      />
      <Divider />
      <p>Linear Integration</p>
      <p>
        Connection Status:{" "}
        {userSettings.linearConnected ? (
          <p className="text-green-600 font-bold"> Connected</p>
        ) : (
          <p className="text-red-600 font-bold"> Not connected</p>
        )}
      </p>
      <Link href={userSettings.linearRedirectURI}>
        <button className="border border-slate-300 text-black py-2 px-4 rounded">
          Connect to linear
        </button>
      </Link>
    </div>
  );
}
