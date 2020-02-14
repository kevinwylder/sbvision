import { session } from './session';

export interface Video {
    id: number
    title: string
    thumbnail: string
    type: number
    format: string
    duration: number
    clips: number
}

interface GetVideosResponse {
    videos: Video[]
    total: number
}

export function getVideos(offset: number, limit: number): Promise<GetVideosResponse> {
    return fetch(`/videos?offset=${offset}&count=${limit}`)
    .then(res => {
        if (res.status != 200) {
            return res.text()
            .then(reason => Promise.reject(reason))
        } 
        return res.json()
    })
    .then(({videos, total}) => videos ? {videos, total} : {videos: [], total: 0});
}

export function getVideoById(id: number) : Promise<GetVideosResponse> {
    return fetch(`/videos?id=${id}`)
    .then(res => res.json())
}

export function addVideo(url: string, type: number) : Promise<Video> {
    return fetch(`/videos`, {
        method: "POST",
        body: JSON.stringify({ type, url }),
        headers: session,
    })
    .then(res => {
        if (res.status != 200) {
            return res.text()
            .then(reason => Promise.reject(reason))
        } 
        return res.json()
    }) 
}