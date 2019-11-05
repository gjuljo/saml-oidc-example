module Main exposing (Model, Msg(..), initialModel, update, view)

import Browser
import Html exposing (..)
import Html.Attributes exposing (action, checked, class, id, method, name, placeholder, src, title, type_, value)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as Decode exposing (Decoder, Value, bool, int, list, string)
import Json.Decode.Pipeline exposing (optional, required)
import Json.Encode as Encode


type alias ResponseFromServer =
    { result : String
    }


type alias Model =
    { content : String
    }


messageToServerEncoder : Model -> Http.Body
messageToServerEncoder model =
    let
        attributes =
            [ ( "content", Encode.string model.content )
            ]
    in
    attributes
        |> Encode.object
        |> Http.jsonBody


responseFromServerDecoder : Decoder ResponseFromServer
responseFromServerDecoder =
    Decode.succeed ResponseFromServer
        |> required "result" string


sendContentToServer : Model -> Cmd Msg
sendContentToServer model =
    let
        body =
            messageToServerEncoder model
    in
    Http.post
        { body = body
        , expect = Http.expectJson HandleResponse responseFromServerDecoder
        , url = "/revert"
        }


initialModel : Model
initialModel =
    { content = ""
    }


view : Model -> Html Msg
view model =
    div
        []
        [ h1 [] [ text "Echo Demo" ]
        , input
            [ type_ "text"
            , placeholder "write something..."
            , onInput ChangeContent
            , value model.content
            ]
            []
        , div [] [ button [ type_ "button", onClick SendContent ] [ text " OK " ] ]
        ]


type Msg
    = ChangeContent String
    | SendContent
    | HandleResponse (Result Http.Error ResponseFromServer)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        ChangeContent newContent ->
            ( { model | content = newContent }, Cmd.none )

        SendContent ->
            ( model, sendContentToServer model )

        HandleResponse (Ok response) ->
            ( { model | content = response.result }, Cmd.none )

        HandleResponse (Err _) ->
            ( model, Cmd.none )


main : Program () Model Msg
main =
    Browser.element
        { init = \_ -> ( initialModel, Cmd.none )
        , subscriptions = \_ -> Sub.none
        , update = update
        , view = view
        }
