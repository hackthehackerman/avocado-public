import Head from "next/head";
import styles from "../styles/Home.module.css";
import Script from "next/script";
import Link from "next/link";
import React, { useState, useEffect } from "react";

export default function Home() {
  return (
    <div className={styles.container}>
      <Head>
        <title>ssssh.....</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Script src="https://accounts.google.com/gsi/client" async defer />

      <main className={styles.main}>
        <div className="flex flex-row max-w-full items-stretch pt-2">
          <p className="font-bold text-2xl text-sky-700">
            Use Linear to manage customer conversations in Slack
          </p>

          <div className=" absolute top-0 right-0">
            <div
              id="g_id_onload"
              data-client_id={process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID}
              data-login_uri={process.env.NEXT_PUBLIC_GOOGLE_OAUTH_REDIRECT_URI}
              data-auto_prompt="false"
            ></div>
            <div
              className="g_id_signin pt-2"
              data-type="standard"
              data-size="large"
              data-theme="outline"
              data-text="sign_in_with"
              data-shape="rectangular"
              data-logo_alignment="left"
            ></div>
          </div>
        </div>

        <iframe
          src="https://www.loom.com/embed/b64449bea50444bca62cc96c5c44851a"
          frameBorder="0"
          allowFullScreen
          className="pt-2"
          width="620"
          height="406"
        ></iframe>

        <div className="pt-2">
          <p className="pt-2 font-bold text-lg">Features and benefits</p>
          <ul className="list-disc">
            <li>
              Each new Slack thread creates an issue in Linear and threaded
              replies within Slack and comments in Linear seamlessly sync both
              ways.
            </li>
            <li>
              {" "}
              React 👀 in Slack to change Linear issue status to &quot;In
              progress.&quot; React ✅ in Slack to change Linear issue status to
              &quot;Done.&quot;
            </li>
            <li>
              Keep your customers updated when their issue status in Linear
              changes with a corresponding emoji that automatically gets added
              to the Slack thread.{" "}
            </li>
            <li>
              Create a searchable repository of past Slack issues in Linear.
            </li>
            <li> Keep track of the status of issue resolutions. </li>
          </ul>
          <p className="pt-2 font-bold text-lg">How does it work?</p>
          <ul className="list-decimal">
            <li>Log-in</li>
            <li>Connect your Slack and Linear accounts</li>
            <li>
              Go to Slack and add the bot “Icarus” to the specific channels you
              wish to integrate with Linear
            </li>
          </ul>
          <p>Have questions?</p>
          <Link href="https://join.slack.com/t/icarus-wgx3901/shared_invite/zt-1fqufq20g-wwDt9qxLBQuWp80SFblqyg">
            <button
              className="flex rounded-md border p-2 justify-around border-gray-300 text-base"
              type="button"
            >
              <div className="flex gap-2 items-center">
                <svg
                  width="20"
                  height="20"
                  viewBox="0 0 54 54"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <g fill="none" fillRule="evenodd">
                    <path
                      d="M19.712.133a5.381 5.381 0 0 0-5.376 5.387 5.381 5.381 0 0 0 5.376 5.386h5.376V5.52A5.381 5.381 0 0 0 19.712.133m0 14.365H5.376A5.381 5.381 0 0 0 0 19.884a5.381 5.381 0 0 0 5.376 5.387h14.336a5.381 5.381 0 0 0 5.376-5.387 5.381 5.381 0 0 0-5.376-5.386"
                      fill="#36C5F0"
                    ></path>
                    <path
                      d="M53.76 19.884a5.381 5.381 0 0 0-5.376-5.386 5.381 5.381 0 0 0-5.376 5.386v5.387h5.376a5.381 5.381 0 0 0 5.376-5.387m-14.336 0V5.52A5.381 5.381 0 0 0 34.048.133a5.381 5.381 0 0 0-5.376 5.387v14.364a5.381 5.381 0 0 0 5.376 5.387 5.381 5.381 0 0 0 5.376-5.387"
                      fill="#2EB67D"
                    ></path>
                    <path
                      d="M34.048 54a5.381 5.381 0 0 0 5.376-5.387 5.381 5.381 0 0 0-5.376-5.386h-5.376v5.386A5.381 5.381 0 0 0 34.048 54m0-14.365h14.336a5.381 5.381 0 0 0 5.376-5.386 5.381 5.381 0 0 0-5.376-5.387H34.048a5.381 5.381 0 0 0-5.376 5.387 5.381 5.381 0 0 0 5.376 5.386"
                      fill="#ECB22E"
                    ></path>
                    <path
                      d="M0 34.249a5.381 5.381 0 0 0 5.376 5.386 5.381 5.381 0 0 0 5.376-5.386v-5.387H5.376A5.381 5.381 0 0 0 0 34.25m14.336-.001v14.364A5.381 5.381 0 0 0 19.712 54a5.381 5.381 0 0 0 5.376-5.387V34.25a5.381 5.381 0 0 0-5.376-5.387 5.381 5.381 0 0 0-5.376 5.387"
                      fill="#E01E5A"
                    ></path>
                  </g>
                </svg>
                <p className="whitespace-nowrap">
                  <b>Join our Slack</b>
                </p>
              </div>
            </button>
          </Link>
        </div>
      </main>
    </div>
  );
}
