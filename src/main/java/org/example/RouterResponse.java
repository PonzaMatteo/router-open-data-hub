package org.example;

public class RouterResponse {

    int statusCode;
    String responseBody;


    public RouterResponse(int statusCode, String responseBody) {
        this.statusCode = statusCode;
        this.responseBody = responseBody;
    }

    public int getStatus() {
        return statusCode;
    }

    public String getBody() {
        return responseBody;
    }
}
