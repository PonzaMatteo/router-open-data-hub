package org.example;

import okhttp3.OkHttpClient;
import okhttp3.Request;

import java.io.IOException;

public class MobilityConnector {

    public static RouterResponse listAllStations() throws IOException, InterruptedException {

        String MOBILITY_URL = "https://mobility.api.opendatahub.com/v2/flat%2Cnode";

        OkHttpClient client = new OkHttpClient();

        Request request = new Request.Builder()
                .url(MOBILITY_URL)
                .build();

        okhttp3.Response response = client.newCall(request).execute();

        // the response:
        RouterResponse routerResponse = new RouterResponse(
                response.code(), response.body().string());

        return routerResponse;
    }


}